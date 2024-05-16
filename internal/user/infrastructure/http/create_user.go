package http

import (
	"context"
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	/* "github.com/hex-api-go/internal/user/application/command/createuser"
	"github.com/hex-api-go/pkg/core/application/cqrs" */)

type Request struct {
	Username string   `validate:"required"`
	Password string   `validate:"required"`
	Devices  []Device `validate:"required,dive,required"`
}
type Device struct {
	Name string `validate:"required"`
}

func validate(req interface{}) {
	var validate = validator.New()
	errs := validate.Struct(req)
	fmt.Println(errs)
}

func CreateUser(ctx context.Context, fiberApp fiber.Router) {
	fiberApp.Post("/create", func(c *fiber.Ctx) error {
		request := new(Request)
		if err := c.BodyParser(request); err != nil {
			return c.SendStatus(400)
		}
		devics :=[]Device{{Name:"iii rapaz"},{}}
		request.Devices = devics
		validate(request)
		/* comand := &createuser.Command{Username: "teste", Password: "123"}
		res, _ := cqrs.Send(comand)
		fmt.Println(res) */
		aa := map[string]any{"teste": "1234"}

		return c.JSON(aa)
	})
}
