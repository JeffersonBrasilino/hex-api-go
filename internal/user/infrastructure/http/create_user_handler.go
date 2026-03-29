package http

import (
	"fmt"

	"github.com/gin-gonic/gin"
	gomes "github.com/jeffersonbrasilino/gomes"
	"github.com/jeffersonbrasilino/gomes/otel"
	"github.com/jeffersonbrasilino/hex-api-go/internal/user/application/command/createuser"
)

var createUserTrace = otel.InitTrace("create-user-handler")

type CreateUserRequest struct {
	Username string `json:"username" binding:"required,gte=18"`
	Password string `json:"password" binding:"required"`
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
			c.Error(err)
			c.AbortWithStatus(400)
			return
		}

		bus, _ := gomes.CommandBus()
		res, err := bus.SendRaw(
			ctx,
			"createUser",
			createuser.CreateCommand(fmt.Sprintf("teste ID: %v", 1), "123"),
			map[string]string{"header1": "val 1", "header2": "val 2"},
		)
		fmt.Println("result", res, err)
		//res, err := bus.Send(ctx, createuser.CreateCommand("teste", "123"))
		if err != nil {
			c.JSON(500, err.Error())
			return
		}

		c.JSON(200, gin.H{"message": "foi OK"})
	})
}
