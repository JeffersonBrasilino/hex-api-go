#### HTTP Handler pattern

HTTP handlers are the entry points for HTTP requests into the module.
They are responsible for request deserialization, input validation, dispatching commands/queries through the `gomes` bus, and returning standardized HTTP responses.

The HTTP handler must:
- be a package-level function that receives a `*gin.RouterGroup` for route registration.
- define a request struct with `json` and `binding` tags for deserialization and validation.
- initialize an OpenTelemetry trace using `gomes/otel.InitTrace`.
- start a span per request with appropriate span kind (`SpanKindServer`).
- dispatch commands via `gomes.CommandBus()` or queries via `gomes.QueryBus()`.
- use `pkg/http` helpers for standardized responses (`http.Error`, `http.ErrorWithCode`, `http.Success`).
- never call domain or repository directly.

Boilerplate Example:

```go
package http

import (
	"fmt"

	httpLib "net/http"

	"github.com/gin-gonic/gin"
	gomes "github.com/jeffersonbrasilino/gomes"
	"github.com/jeffersonbrasilino/gomes/otel"
	"github.com/jeffersonbrasilino/hex-api-go/internal/[module-name]/application/command/[actionname]"
	"github.com/jeffersonbrasilino/hex-api-go/pkg/http"
)

var [actionName]Trace = otel.InitTrace("[action-name]-handler")

type [ActionName]Request struct {
	Field1 string `json:"field1" binding:"required"`
	Field2 string `json:"field2" binding:"required"`
}

func [ActionName]Handler(router *gin.RouterGroup) {
	uri := "/[action-uri]"
	router.POST(uri, func(c *gin.Context) {
		ctx, span := [actionName]Trace.Start(
			c,
			fmt.Sprintf("post %s", uri),
			otel.WithSpanKind(otel.SpanKindServer),
		)
		defer span.End()

		var request [ActionName]Request
		if err := c.ShouldBindJSON(&request); err != nil {
			http.ErrorWithCode(c, httpLib.StatusBadRequest, err)
			return
		}

		bus, _ := gomes.CommandBus()
		res, err := bus.Send(ctx, &[actionname].Command{
			Field1: request.Field1,
			Field2: request.Field2,
		})

		if err != nil {
			http.Error(c, err)
			return
		}

		http.Success(c, httpLib.StatusOK, res)
	})
}
```
Implementation example: see -> `../../internal/user/infrastructure/http/create_user_handler.go`
