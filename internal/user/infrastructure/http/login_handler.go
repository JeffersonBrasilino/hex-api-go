// Login HTTP handler.
//
// Intent: expose the login use case over HTTP.
// Objective: bind the submitted credentials, dispatch the login command through the bus, and
// return a standardized response — 429 when the user is locked out, or the generic error mapping
// for any other failure.
package http

import (
	"fmt"

	httpLib "net/http"

	"github.com/gin-gonic/gin"
	gomes "github.com/jeffersonbrasilino/gomes"
	"github.com/jeffersonbrasilino/gomes/otel"
	"github.com/jeffersonbrasilino/hex-api-go/internal/user/application/command/login"
	"github.com/jeffersonbrasilino/hex-api-go/internal/user/domain"
	"github.com/jeffersonbrasilino/hex-api-go/pkg/http"
)

var loginTrace = otel.InitTrace("login-handler")

// LoginRequest is the payload expected by LoginHandler.
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginHandler registers the login route on the given router group.
//
// Intent: authenticate a user by dispatching a login.Command through the command bus.
// Parameters:
//   - router: the gin.RouterGroup to register the route on.
//
// Behavior: binds the request JSON, returns 400 on binding failure, 429 when the returned error
// is a *domain.UserBlockedError, the generic mapped status for any other error, or 201 with the
// issued session data on success.
func LoginHandler(router *gin.RouterGroup) {
	uri := "/login"
	router.POST(uri, func(c *gin.Context) {
		ctx, span := loginTrace.Start(
			c,
			fmt.Sprintf("post %s", uri),
			otel.WithSpanKind(otel.SpanKindServer),
		)
		defer span.End()

		var request LoginRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			http.ErrorWithCode(c, httpLib.StatusBadRequest, err)
			return
		}

		bus, _ := gomes.CommandBus()
		res, err := bus.Send(ctx, &login.Command{
			Username: request.Username,
			Password: request.Password,
		})

		if err != nil {
			if _, ok := err.(*domain.UserBlockedError); ok {
				http.ErrorWithCode(c, httpLib.StatusTooManyRequests, err)
				return
			}
			http.Error(c, err)
			return
		}

		http.Success(c, httpLib.StatusCreated, res)
	})
}
