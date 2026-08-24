// Revoke session HTTP handler.
//
// Intent: expose the administrative session revocation use case over HTTP.
// Objective: read the session id from the URL, dispatch the revokeSession command through the
// bus, and return a standardized response.
package http

import (
	"fmt"

	httpLib "net/http"

	"github.com/gin-gonic/gin"
	gomes "github.com/jeffersonbrasilino/gomes"
	"github.com/jeffersonbrasilino/gomes/otel"
	"github.com/jeffersonbrasilino/hex-api-go/internal/user/application/command/revokesession"
	"github.com/jeffersonbrasilino/hex-api-go/pkg/http"
)

var revokeSessionTrace = otel.InitTrace("revoke-session-handler")

// RevokeSessionHandler registers the administrative session revocation route on the given
// router group.
//
// Intent: end any user's session by dispatching a revokesession.Command through the command bus.
// Parameters:
//   - router: the gin.RouterGroup to register the route on.
//
// Behavior: reads the sessionId path parameter, dispatches the command, returns the generic
// mapped status (404 for a non-existent session) on error, or 200 on success.
func RevokeSessionHandler(router *gin.RouterGroup) {
	uri := "/logout/:sessionId"
	router.DELETE(uri, func(c *gin.Context) {
		ctx, span := revokeSessionTrace.Start(
			c,
			fmt.Sprintf("delete %s", uri),
			otel.WithSpanKind(otel.SpanKindServer),
		)
		defer span.End()

		sessionId := c.Param("sessionId")

		bus, _ := gomes.CommandBus()
		res, err := bus.Send(ctx, &revokesession.Command{
			SessionId: sessionId,
		})

		if err != nil {
			http.Error(c, err)
			return
		}

		http.Success(c, httpLib.StatusOK, res)
	})
}
