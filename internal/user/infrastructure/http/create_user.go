package http

import (
	"context"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/hex-api-go/internal/user/application/command/createuser"
	"github.com/hex-api-go/internal/user/application/query/getuser"
	messagesystem "github.com/hex-api-go/pkg/core/infrastructure/messagesystem"
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

		//coreHttp.ValidateRequest(request)

		/* bus := messagesystem.CommandBus()
		res, err := bus.Send(c.Context(), createuser.CreateCommand("teste", "123")) */

		//opCtx, cancel := context.WithTimeout(c.Context(), time.Second*5)
		//defer cancel()

		busA := messagesystem.CommandBusByChannel("messagesystem.topic")
		err := busA.SendAsync(c.Context(), createuser.CreateCommand("teste", "123"))
		fmt.Println("[controller] ASYNC COMMAND SEND ERROR ", err)

		//ch := map[string]string{"foi": "okokokokok"}
		//err := busA.SendRawAsync(opCtx, "createUser", "iii rapaz DEU BOM", ch)

		/* if err != nil {
			fmt.Println(time.Now().Format("2006-01-02 15:04:05"), "[controller] ERROR ", err)
			return c.SendStatus(500)
		} */
		return c.JSON("foi OK")
	})

	fiberApp.Get("", func(c *fiber.Ctx) error {
		/* request := new(Request)
		if err := c.BodyParser(request); err != nil {
			return c.SendStatus(400)
		} */

		//coreHttp.ValidateRequest(request)
		bus := messagesystem.QueryBus()
		res, err := bus.Send(c.Context(), getuser.NewQuery())
		//a, ok := res.(*createuser.ResultCm)
		fmt.Println("controller", res, err)

		return c.JSON("okokokok")
	})
}
