package http

import (
	"context"
	"encoding/json"

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

		command, _ := json.Marshal(createuser.CreateCommand("teste", "123"))
		commandBus := messagesystem.GetCommandBus()
		commandBus.Send("CreateUser", command, nil)

		coreHttp.ValidateRequest(request)

		/* comand := createuser.CreateCommand("teste", "123")
		res, err := cqrs.Send(comand)
		if err != nil {
			var message any
			json.Unmarshal([]byte(err.Error()), &message)
			return c.Status(400).JSON(message)
		} */

		return c.JSON("okokokokko")
	})
}
