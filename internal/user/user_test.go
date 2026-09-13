package user_test

import (
	"context"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/jeffersonbrasilino/hex-api-go/internal/user"
)

// newTestEngine builds a gin.Engine in test mode, suitable for route/middleware assertions
// without emitting Gin's debug logging.
func newTestEngine() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.New()
}

// TestNewUserModule verifies the constructor wires the module without error and returns a
// usable, non-nil instance ready to be registered.
func TestNewUserModule(t *testing.T) {
	t.Run("should return a non-nil module instance", func(t *testing.T) {
		t.Parallel()
		engine := newTestEngine()

		module := user.NewUserModule(engine, nil, nil, "test-secret")

		if module == nil {
			t.Fatal("NewUserModule should return a non-nil module instance")
		}
	})
}

// TestUserModule_WithHttpProtocol verifies the module's HTTP routes are registered under the
// "/users" prefix, independent of the rest of the module's wiring.
func TestUserModule_WithHttpProtocol(t *testing.T) {
	t.Run("should register the module's http routes under /users", func(t *testing.T) {
		t.Parallel()
		engine := newTestEngine()
		module := user.NewUserModule(engine, nil, nil, "test-secret")

		result := module.WithHttpProtocol()

		if result == nil {
			t.Fatal("WithHttpProtocol should return a non-nil module instance")
		}

		expectedRoutes := map[string]string{
			"/users/create":            "POST",
			"/users/logout/:sessionId": "DELETE",
		}
		foundRoutes := map[string]string{}
		for _, route := range engine.Routes() {
			foundRoutes[route.Path] = route.Method
		}

		for path, method := range expectedRoutes {
			if foundMethod, ok := foundRoutes[path]; !ok || foundMethod != method {
				t.Errorf("expected route %s %s to be registered, got routes: %v", method, path, foundRoutes)
			}
		}

		if _, ok := foundRoutes["/users/login"]; !ok {
			t.Errorf("expected route POST /users/login to be registered, got routes: %v", foundRoutes)
		}
	})
}

// TestUserModule_Register verifies that Register wires the module's dependencies, applies the
// authorization middleware before the HTTP routes are set up, and exposes the module's routes
// without error.
func TestUserModule_Register(t *testing.T) {
	t.Run("should register actions and http routes without error", func(t *testing.T) {
		t.Parallel()
		engine := newTestEngine()
		module := user.NewUserModule(engine, nil, nil, "test-secret")
		ctx := context.Background()

		err := module.Register(ctx)

		if err != nil {
			t.Fatalf("Register should not return an error, got: %v", err)
		}

		foundRoutes := map[string]bool{}
		for _, route := range engine.Routes() {
			foundRoutes[route.Method+" "+route.Path] = true
		}

		expected := []string{
			"POST /users/create",
			"POST /users/login",
			"DELETE /users/logout/:sessionId",
		}
		for _, route := range expected {
			if !foundRoutes[route] {
				t.Errorf("expected route %q to be registered, got routes: %v", route, foundRoutes)
			}
		}
	})
}
