package http

import (
	"fmt"

	httpLib "net/http"

	"github.com/gin-gonic/gin"
	gomes "github.com/jeffersonbrasilino/gomes"
	"github.com/jeffersonbrasilino/gomes/otel"
	"github.com/jeffersonbrasilino/hex-api-go/internal/user/application/command/createuser"
	"github.com/jeffersonbrasilino/hex-api-go/pkg/http"
)

var createUserTrace = otel.InitTrace("create-user-handler")

type CreateUserRequest struct {
	Username   string `json:"username" binding:"required"`
	Password   string `json:"password" binding:"required"`
	PersonName string `json:"name" binding:"required"`
	Document   string `json:"document" binding:"required"`
	BirthDate  string `json:"birthDate" binding:"required"`
	Email      string `json:"email" binding:"required"`
}

func CreateUserHandler(router *gin.RouterGroup) {
	uri := "/create"
	router.POST(uri, func(c *gin.Context) {
		ctx, span := createUserTrace.Start(
			c,
			fmt.Sprintf("post %s", uri),
			otel.WithSpanKind(otel.SpanKindServer),
		)
		defer span.End()

		var request CreateUserRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			http.ErrorWithCode(c, httpLib.StatusBadRequest, err)
			return
		}

		bus, _ := gomes.CommandBus()
		res, err := bus.Send(ctx, &createuser.Command{
			Username:   request.Username,
			Password:   request.Password,
			PersonName: request.PersonName,
			Document:   request.Document,
			BirthDate:  request.BirthDate,
			Email:      request.Email,
		})

		if err != nil {
			http.Error(c, err)
			return
		}

		http.Success(c, httpLib.StatusOK, res)
	})
}
