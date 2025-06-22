package http

import (
	"context"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/hex-api-go/internal/user/application/command/createuser"
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

		/* time.Sleep(3 * time.Second)
		fmt.Println(time.Now().Format("2006-01-02 15:04:05"), "[controller] OK ")
		return c.JSON("okokokokokokokok") */

		//opCtx, cancel := context.WithTimeout(c.Context(), time.Second*4)
		//defer cancel()
		//coreHttp.ValidateRequest(request)

		//COMMAND SYNC
		bus := messagesystem.CommandBus()
		res, err := bus.Send(c.Context(), createuser.CreateCommand("teste", "123"))
		//COMMAND ASYNC
		/* 	opCtx, cancel := context.WithTimeout(c.Context(), time.Second*5)
		defer cancel()
		bus := messagesystem.CommandBusByChannel("message_system.topic")
		err := bus.SendAsync(opCtx, createuser.CreateCommand("teste", "123"))
		fmt.Println("[controller] ASYNC COMMAND SEND ERROR ", err) */

		//RAW COMMAND ASYNC
		/* createuser.CreateCommand("teste", "123")
		bus := messagesystem.CommandBusByChannel("message_system.topic")
		err := bus.SendRawAsync(opCtx, "raw_message_route", "huehuebrbr", nil)
		fmt.Println("[controller] ASYNC COMMAND SEND ERROR ", err) */

		/* evBus := messagesystem.EventBusByChannel("message_system.topic")
		erra := evBus.Publish(createuser.NewCreatedCommand("teste", "123"))
		fmt.Println("[event-bus] publish error ", erra) */

		if err != nil {
			fmt.Println(time.Now().Format("2006-01-02 15:04:05"), "[controller] ERROR ", err)
			return c.SendStatus(500)
		}

		fmt.Println(time.Now().Format("2006-01-02 15:04:05"), "[controller] OK ")
		return c.JSON(res)
	})

	fiberApp.Get("", func(c *fiber.Ctx) error {
		/* request := new(Request)
		if err := c.BodyParser(request); err != nil {
			return c.SendStatus(400)
		} */

		//coreHttp.ValidateRequest(request)
		/* 		bus := messagesystem.QueryBus()
		   		res, err := bus.Send(getuser.NewQuery())
		   		//a, ok := res.(*createuser.ResultCm)
		   		fmt.Println("controller", res, err) */

		return c.JSON("okokokok")
	})
}
