// Login HTTP handler tests.
//
// Intent: verify the HTTP boundary for the login route in isolation from the real login use case.
// Objective: cover request binding failure, the 429 lockout mapping, the generic error mapping,
// and the successful session-issuance response.
//
// Bus bootstrap: gomes.CommandBus() is a process-wide singleton that can only be started once
// (github.com/jeffersonbrasilino/gomes/gomes.go — Start() errors if called a second time), so this
// file's TestMain bootstraps it once for the whole http_test package, registering stub action
// handlers for login, logout, and revokeSession — shared by login_handler_test.go,
// logout_handler_test.go, and revoke_session_handler_test.go. Each stub's behavior is a pure
// function of the command payload (no shared mutable state), so subtests across all three files
// remain safe to run with t.Parallel().
package http_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/jeffersonbrasilino/ddgo"
	"github.com/jeffersonbrasilino/hex-api-go/internal/user/application/command/login"
	"github.com/jeffersonbrasilino/hex-api-go/internal/user/application/command/revokesession"
	"github.com/jeffersonbrasilino/hex-api-go/internal/user/domain"
	httpHandler "github.com/jeffersonbrasilino/hex-api-go/internal/user/infrastructure/http"
)

// blockedLoginUsername, invalidLoginUsername select the stubLoginActionHandler's deterministic
// outcome by username, so tests can run in parallel without mutating shared state.
const (
	blockedLoginUsername = "blocked-user"
	invalidLoginUsername = "invalid-user"
)

// stubLoginActionHandler is a hand-rolled test double registered on the real command bus,
// standing in for application/command/login.Handler so this file exercises only the HTTP layer.
type stubLoginActionHandler struct{}

func (h *stubLoginActionHandler) Handle(ctx context.Context, cmd *login.Command) (any, error) {
	switch cmd.Username {
	case blockedLoginUsername:
		return nil, domain.NewUserBlockedError("Usuário temporariamente bloqueado por excesso de tentativas")
	case invalidLoginUsername:
		return nil, ddgo.NewValidationError("Credenciais inválidas")
	default:
		return map[string]string{
			"accessToken":  "access-token-value",
			"refreshToken": "refresh-token-value",
			"sessionId":    "session-id-value",
		}, nil
	}
}

// notFoundSessionId selects the stubRevokeSessionActionHandler's deterministic outcome.
const notFoundSessionId = "not-found-session"

// stubRevokeSessionActionHandler is a hand-rolled test double registered on the real command bus,
// standing in for application/command/revokesession.Handler so revoke_session_handler_test.go
// exercises only the HTTP layer.
type stubRevokeSessionActionHandler struct{}

func (h *stubRevokeSessionActionHandler) Handle(ctx context.Context, cmd *revokesession.Command) (any, error) {
	if cmd.SessionId == notFoundSessionId {
		return nil, ddgo.NewNotFoundError("Sessão não encontrada")
	}
	return "session revoked", nil
}

// newLoginRouter builds a gin engine with the login route registered, ready to receive test
// requests.
func newLoginRouter() *gin.Engine {
	router := gin.New()
	group := router.Group("")
	httpHandler.LoginHandler(group)
	return router
}

func TestLoginHandler(t *testing.T) {
	t.Run("should return 400 when the request body fails to bind", func(t *testing.T) {
		t.Parallel()
		router := newLoginRouter()
		req := httptest.NewRequest("POST", "/login", bytes.NewBufferString(`{"username":""}`))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != 400 {
			t.Fatalf("expected status 400 for a binding failure, got: %d, body: %s", w.Code, w.Body.String())
		}
	})

	t.Run("should return 429 when the login use case reports the user is blocked", func(t *testing.T) {
		t.Parallel()
		router := newLoginRouter()
		body, _ := json.Marshal(map[string]string{"username": blockedLoginUsername, "password": "any-password"})
		req := httptest.NewRequest("POST", "/login", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != 429 {
			t.Fatalf("expected status 429 when the user is blocked, got: %d, body: %s", w.Code, w.Body.String())
		}
	})

	t.Run("should map non-blocked errors through the generic error mapping", func(t *testing.T) {
		t.Parallel()
		router := newLoginRouter()
		body, _ := json.Marshal(map[string]string{"username": invalidLoginUsername, "password": "any-password"})
		req := httptest.NewRequest("POST", "/login", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != 400 {
			t.Fatalf("expected status 400 for a ddgo.ValidationError, got: %d, body: %s", w.Code, w.Body.String())
		}
	})

	t.Run("should return 201 with the issued session on success", func(t *testing.T) {
		t.Parallel()
		router := newLoginRouter()
		body, _ := json.Marshal(map[string]string{"username": "jdoe", "password": "StrongP@ss1"})
		req := httptest.NewRequest("POST", "/login", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != 201 {
			t.Fatalf("expected status 201 on success, got: %d, body: %s", w.Code, w.Body.String())
		}
		var payload map[string]string
		if err := json.Unmarshal(w.Body.Bytes(), &payload); err != nil {
			t.Fatalf("response body should be valid JSON, got: %v", err)
		}
		if payload["accessToken"] != "access-token-value" {
			t.Fatalf("expected accessToken in the response, got: %v", payload)
		}
		if payload["refreshToken"] != "refresh-token-value" {
			t.Fatalf("expected refreshToken in the response, got: %v", payload)
		}
		if payload["sessionId"] != "session-id-value" {
			t.Fatalf("expected sessionId in the response, got: %v", payload)
		}
	})
}
