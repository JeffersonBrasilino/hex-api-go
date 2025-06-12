package http

import (
	"context"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/hex-api-go/internal/user/application/command/createuser"
	"github.com/hex-api-go/internal/user/application/query/getuser"
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

		//COMMAND SYNC
		bus := messagesystem.CommandBus()
		res, err := bus.Send(createuser.CreateCommand("teste", "123"))
		fmt.Println("[controller] result command ", res, err)
		
		//COMMAND ASYNC
		/* bus := messagesystem.CommandBusByChannel("message_system.topic")
		err := bus.SendAsync(createuser.CreateCommand("teste", "123"))
		fmt.Println("[controller] ASYNC COMMAND SEND ERROR ", err) */

		//RAW COMMAND ASYNC
		/* createuser.CreateCommand("teste", "123")
		bus := messagesystem.CommandBusByChannel("message_system.topic")
		err := bus.SendRawAsync("raw_message_route","huehuebrbr", nil)
		fmt.Println("[controller] ASYNC COMMAND SEND ERROR ", err) */
		
		/* evBus := messagesystem.EventBusByChannel("message_system.topic")
		erra := evBus.Publish(createuser.NewCreatedCommand("teste", "123"))
		fmt.Println("[event-bus] publish error ", erra) */
		
		return c.JSON("okokokokokokokok")
	})

	fiberApp.Get("", func(c *fiber.Ctx) error {
		/* request := new(Request)
		if err := c.BodyParser(request); err != nil {
			return c.SendStatus(400)
		} */

		//coreHttp.ValidateRequest(request)
		bus := messagesystem.QueryBus()
		res, err := bus.Send(getuser.NewQuery())
		//a, ok := res.(*createuser.ResultCm)
		fmt.Println("controller", res, err)

		return c.JSON("okokokok")
	})
}
