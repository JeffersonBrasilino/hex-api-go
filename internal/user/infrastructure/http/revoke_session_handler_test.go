// Revoke session HTTP handler tests.
//
// Intent: verify the HTTP boundary for the administrative session revocation route in isolation
// from the real revokeSession use case.
// Objective: cover the 404 not-found mapping and the successful revocation response.
//
// The shared command bus and the stubRevokeSessionActionHandler registered on it are bootstrapped
// once by TestMain in login_handler_test.go — see that file's package doc for why.
package http_test

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	httpHandler "github.com/jeffersonbrasilino/hex-api-go/internal/user/infrastructure/http"
)

// newRevokeSessionRouter builds a gin engine with the revoke session route registered, ready to
// receive test requests.
func newRevokeSessionRouter() *gin.Engine {
	router := gin.New()
	group := router.Group("")
	httpHandler.RevokeSessionHandler(group)
	return router
}

func TestRevokeSessionHandler(t *testing.T) {
	t.Run("should return 404 when the revokeSession use case reports the session was not found", func(t *testing.T) {
		t.Parallel()
		router := newRevokeSessionRouter()
		req := httptest.NewRequest("DELETE", "/logout/"+notFoundSessionId, nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != 404 {
			t.Fatalf("expected status 404 for a non-existent session, got: %d, body: %s", w.Code, w.Body.String())
		}
	})

	t.Run("should return 200 on successful session revocation", func(t *testing.T) {
		t.Parallel()
		router := newRevokeSessionRouter()
		req := httptest.NewRequest("DELETE", "/logout/an-existing-session", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != 200 {
			t.Fatalf("expected status 200 on success, got: %d, body: %s", w.Code, w.Body.String())
		}
	})
}
