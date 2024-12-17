package http

import (
	"context"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/hex-api-go/internal/user/application/command/createuser"
	coreHttp "github.com/hex-api-go/pkg/core/infrastructure/http"
	messagesystem "github.com/hex-api-go/pkg/core/infrastructure/message_system"
)

type Request struct {
	Username string `validate:"gte=4"`
	Password string `validate:"required"`
}

func CreateUser(ctx context.Context, fiberApp fiber.Router) {
	fiberApp.Post("/create", func(c *fiber.Ctx) error {
		request := new(Request)
		if err := c.BodyParser(request); err != nil {
			return c.SendStatus(400)
		}

		coreHttp.ValidateRequest(request)
		bus := messagesystem.GetCommandBus()
		res, _ := bus.Send(createuser.CreateCommand("teste", "123"))
		a, ok := res.(*createuser.ResultCm)
		fmt.Println("controller", a, ok, res)

		return c.JSON(res)
	})
}
