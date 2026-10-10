// Authorization HTTP middleware.
//
// Intent: enforce RBAC access control at the HTTP layer by extracting the
// access token from the Authorization header, checking user permissions against
// the requested resource (HTTP method + path), and allowing or denying access
// based on the permission check result.
// Objective: provide a middleware, registered on protected router groups/routes,
// that validates access tokens and dispatches permission checks through the
// command bus (never calling domain or repository directly), mapping errors to
// standardized HTTP status codes (401 for invalid tokens, 403 for access
// denied, others via generic error mapping).
package http

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jeffersonbrasilino/gomes"
	"github.com/jeffersonbrasilino/gomes/otel"
	"github.com/jeffersonbrasilino/hex-api-go/internal/user/application/query/checkpermission"
	"github.com/jeffersonbrasilino/hex-api-go/internal/user/domain"
	httpLib "github.com/jeffersonbrasilino/hex-api-go/pkg/http"
)

var authorizationTrace = otel.InitTrace("authorization-middleware")

// AuthorizationMiddleware returns a Gin middleware that enforces RBAC access
// control by checking user permissions for each HTTP request.
//
// Intent: create a middleware that checks permissions through the CQRS bus for every
// request, denying access on token failure or insufficient permissions, without ever
// deciding on its own whether a route requires a token — that decision belongs entirely
// to the checkpermission handler, which is the only place that knows whether a route has
// any permission mapped at all (RF-04's data-driven public-route rule).
//
// Returns: a gin.HandlerFunc that:
//   - Extracts the bearer token from the Authorization header if present and
//     well-formed; otherwise dispatches with an empty access token instead of rejecting
//     the request outright — a missing token is only an error if the target route turns
//     out to require one, which only the handler can determine.
//   - Dispatches a checkpermission.Query through the command bus for every request.
//   - Returns 401 if the route is protected and the token is missing/invalid/expired
//     (InvalidSessionError).
//   - Returns 403 if the user lacks permissions (AccessDeniedError).
//   - Returns the generic error mapping for any other infrastructure failure.
//   - Calls c.Next() and continues the chain on success (access allowed or public
//     route).
//   - Aborts the chain (c.Abort()) on every error branch, so an infrastructure
//     failure never falls through to the protected route handler.
//
// Behavior: This middleware must be registered explicitly on the router
// groups/routes that require access control — never globally — to avoid
// checking permissions for public routes (e.g. login) that need no such
// overhead.
func AuthorizationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := authorizationTrace.Start(
			c,
			"authorization-check",
			otel.WithSpanKind(otel.SpanKindServer),
		)
		defer span.End()

		// Extract the bearer token if the Authorization header is present and well-formed;
		// otherwise proceed with an empty token. Whether that is actually a problem is for
		// the handler to decide, once it knows if the route is protected.
		token := ""
		if authHeader := c.GetHeader("Authorization"); authHeader != "" {
			if parts := strings.SplitN(authHeader, " ", 2); len(parts) == 2 &&
				strings.ToLower(parts[0]) == "bearer" {
				token = parts[1]
			}
		}

		// Dispatch checkpermission query through the bus.
		bus, _ := gomes.QueryBus()
		_, err := bus.Send(ctx, &checkpermission.Query{
			AccessToken: token,
			Method:      c.Request.Method,
			Path:        c.FullPath(),
		})

		if err != nil {
			// Map domain errors to HTTP status codes. c.Abort() is required in every
			// branch: Gin's handler chain is driven by an outer for-loop in Next() that
			// keeps invoking subsequent handlers once this one returns, regardless of
			// whether c.Next() was called here — only c.Abort() stops it.
			if _, ok := err.(*domain.InvalidSessionError); ok {
				httpLib.ErrorWithCode(c, http.StatusUnauthorized, err)
				c.Abort()
				return
			}
			if _, ok := err.(*domain.AccessDeniedError); ok {
				httpLib.ErrorWithCode(c, http.StatusForbidden, err)
				c.Abort()
				return
			}
			// All other errors (infrastructure failures) use generic mapping.
			c.Error(err)
			c.Abort()
			return
		}

		// Access allowed (user has permission or route is public).
		c.Next()
	}
}
