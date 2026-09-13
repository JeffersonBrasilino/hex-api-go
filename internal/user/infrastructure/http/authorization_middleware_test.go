// Authorization middleware tests.
//
// Intent: verify the HTTP boundary for the authorization middleware in isolation
// from the real permission checking use case.
// Objective: cover Authorization header extraction, Bearer token parsing, error
// mapping (401 for invalid tokens, 403 for access denied, others via generic
// mapping), successful access granted (c.Next() called), and — critically — that a
// missing/malformed header never short-circuits the request before the bus is asked
// whether the route is even protected (a public route must remain reachable by an
// unauthenticated caller).
//
// Bus bootstrap: gomes.QueryBus() is backed by a process-wide singleton that
// can only be started once. The http_test package's TestMain (from
// login_handler_test.go) bootstraps the bus once, and this file registers a
// stub query handler via init() for checkpermission.Query. The stub's behavior
// is a pure function of the query payload (token, method, path), so subtests
// can run with t.Parallel().
package http_test

import (
	"context"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	gomes "github.com/jeffersonbrasilino/gomes"
	"github.com/jeffersonbrasilino/hex-api-go/internal/user/application/query/checkpermission"
	"github.com/jeffersonbrasilino/hex-api-go/internal/user/domain"
	httpHandler "github.com/jeffersonbrasilino/hex-api-go/internal/user/infrastructure/http"
)

// init registers the stub checkpermission query handler once, before any test
// runs. The handler will be available when the bus is started by login_handler_test.go's
// TestMain.
func init() {
	if err := gomes.AddActionHandler[*checkpermission.Query, any](&stubCheckPermissionHandler{}); err != nil {
		// Handler registration may fail if called multiple times or if the bus is
		// already started; the login_handler_test.go TestMain ensures the bus
		// starts after all handlers are registered.
		_ = err
	}
}

// stubCheckPermissionHandler is a hand-rolled test double registered on the real
// query bus, standing in for the checkpermission.Handler so this file exercises
// only the HTTP layer.
type stubCheckPermissionHandler struct{}

func (h *stubCheckPermissionHandler) Handle(
	ctx context.Context,
	q *checkpermission.Query,
) (any, error) {
	// "/public" simulates a route with no permission mapped — allowed regardless of
	// whether a token was presented at all, mirroring the real handler's RN-01 order
	// (route mapping is checked before the token).
	if q.Path == "/public" {
		return true, nil
	}

	// Every other path is a stand-in for a protected route. Behavior is a deterministic
	// function of the token value, allowing tests to run in parallel without shared
	// state.
	switch q.AccessToken {
	case "":
		// No token presented for a route that requires one.
		return nil, domain.NewInvalidSessionError("Missing access token")
	case "invalid-token":
		return nil, domain.NewInvalidSessionError("Invalid or expired token")
	case "expired-token":
		return nil, domain.NewInvalidSessionError("Token expired")
	case "denied-token":
		return nil, domain.NewAccessDeniedError("User does not have permission to access this resource")
	case "infrastructure-error":
		return nil, fmt.Errorf("database connection failed")
	default:
		// Default: successful authorization (user has permission)
		return true, nil
	}
}

// newAuthRouter builds a gin engine with a protected route and a public route for
// testing the authorization middleware.
func newAuthRouter() *gin.Engine {
	router := gin.New()

	// Apply authorization middleware globally.
	router.Use(httpHandler.AuthorizationMiddleware())

	// Protected route for testing.
	router.GET("/protected", func(c *gin.Context) {
		c.JSON(200, map[string]string{"message": "access granted"})
	})

	// Public route (no permission mapped) for testing unauthenticated access.
	router.GET("/public", func(c *gin.Context) {
		c.JSON(200, map[string]string{"message": "public access granted"})
	})

	return router
}

func TestAuthorizationMiddleware(t *testing.T) {
	t.Run("missing Authorization header on a protected route returns 401 via the bus", func(t *testing.T) {
		t.Parallel()

		router := newAuthRouter()
		req := httptest.NewRequest("GET", "/protected", nil)
		// No Authorization header set. The middleware must still dispatch the query
		// (with an empty AccessToken) instead of short-circuiting — 401 here comes from
		// the stub handler treating an empty token on a protected route as invalid, not
		// from the middleware rejecting the request on its own.
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != 401 {
			t.Fatalf("expected status 401 for missing header on a protected route, got: %d, body: %s",
				w.Code, w.Body.String())
		}
	})

	t.Run("malformed Authorization header on a protected route returns 401 via the bus", func(t *testing.T) {
		t.Parallel()

		router := newAuthRouter()
		req := httptest.NewRequest("GET", "/protected", nil)
		// Malformed header: missing Bearer prefix — the middleware treats this the same
		// as no token at all (empty AccessToken dispatched), not as an immediate reject.
		req.Header.Set("Authorization", "malformed-token-without-bearer-prefix")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != 401 {
			t.Fatalf("expected status 401 for malformed header, got: %d, body: %s",
				w.Code, w.Body.String())
		}
	})

	t.Run("invalid Bearer format on a protected route returns 401 via the bus", func(t *testing.T) {
		t.Parallel()

		router := newAuthRouter()
		req := httptest.NewRequest("GET", "/protected", nil)
		// Invalid format: only "Bearer" without token — treated as empty AccessToken.
		req.Header.Set("Authorization", "Bearer")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != 401 {
			t.Fatalf("expected status 401 for invalid Bearer format, got: %d, body: %s",
				w.Code, w.Body.String())
		}
	})

	t.Run("invalid or expired token returns 401", func(t *testing.T) {
		t.Parallel()

		router := newAuthRouter()
		req := httptest.NewRequest("GET", "/protected", nil)
		req.Header.Set("Authorization", "Bearer invalid-token")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != 401 {
			t.Fatalf("expected status 401 for invalid token, got: %d, body: %s",
				w.Code, w.Body.String())
		}
	})

	t.Run("expired token returns 401", func(t *testing.T) {
		t.Parallel()

		router := newAuthRouter()
		req := httptest.NewRequest("GET", "/protected", nil)
		req.Header.Set("Authorization", "Bearer expired-token")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != 401 {
			t.Fatalf("expected status 401 for expired token, got: %d, body: %s",
				w.Code, w.Body.String())
		}
	})

	t.Run("access denied returns 403", func(t *testing.T) {
		t.Parallel()

		router := newAuthRouter()
		req := httptest.NewRequest("GET", "/protected", nil)
		req.Header.Set("Authorization", "Bearer denied-token")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != 403 {
			t.Fatalf("expected status 403 for access denied, got: %d, body: %s",
				w.Code, w.Body.String())
		}
	})

	t.Run("infrastructure error returns generic error mapping", func(t *testing.T) {
		t.Parallel()

		router := newAuthRouter()
		req := httptest.NewRequest("GET", "/protected", nil)
		req.Header.Set("Authorization", "Bearer infrastructure-error")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		// Infrastructure errors map to 500 via the generic error mapping.
		if w.Code != 500 {
			t.Fatalf("expected status 500 for infrastructure error, got: %d, body: %s",
				w.Code, w.Body.String())
		}
	})

	t.Run("valid token with permission calls c.Next() and proceeds", func(t *testing.T) {
		t.Parallel()

		router := newAuthRouter()
		req := httptest.NewRequest("GET", "/protected", nil)
		req.Header.Set("Authorization", "Bearer valid-token")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != 200 {
			t.Fatalf("expected status 200 for valid token, got: %d, body: %s",
				w.Code, w.Body.String())
		}
	})

	t.Run("public route (no permission mappings) allows access with a token", func(t *testing.T) {
		t.Parallel()

		router := newAuthRouter()
		req := httptest.NewRequest("GET", "/public", nil)
		req.Header.Set("Authorization", "Bearer any-token")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != 200 {
			t.Fatalf("expected status 200 for public route, got: %d, body: %s",
				w.Code, w.Body.String())
		}
	})

	t.Run("public route (no permission mappings) allows access with no token at all", func(t *testing.T) {
		t.Parallel()

		router := newAuthRouter()
		req := httptest.NewRequest("GET", "/public", nil)
		// No Authorization header set — the middleware must not reject the request just
		// for lacking a token; only the handler, which knows the route is unmapped
		// (public), gets to decide. This is the regression this test guards against: an
		// unauthenticated call to a public route (e.g. login) must never be rejected for
		// missing a token it was never meant to carry.
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != 200 {
			t.Fatalf("expected status 200 for an unauthenticated request to a public route, got: %d, body: %s",
				w.Code, w.Body.String())
		}
	})

	t.Run("case-insensitive Bearer prefix", func(t *testing.T) {
		t.Parallel()

		router := newAuthRouter()
		req := httptest.NewRequest("GET", "/protected", nil)
		req.Header.Set("Authorization", "bearer valid-token")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != 200 {
			t.Fatalf("expected status 200 for case-insensitive Bearer, got: %d, body: %s",
				w.Code, w.Body.String())
		}
	})
}
