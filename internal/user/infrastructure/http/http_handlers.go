package http

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/hex-api-go/internal/user/application/command/createuser"
	"github.com/hex-api-go/internal/user/application/query/getuser"
	"github.com/hex-api-go/pkg/core/application/cqrs"
)

func NewHttpHandlers(app *fiber.App) {
	router := app.Group("/users")
	router.Post("/create", createUser)
	router.Get("", getUser)
}

func createUser(c *fiber.Ctx) error {
	comand := &createuser.Command{Username: "teste", Password: "123"}
	res, _ := cqrs.Send(comand)

	fmt.Println(res)

	return c.JSON(res)
}

func getUser(c *fiber.Ctx) error {
	res, err := cqrs.Send(&getuser.Query{DataSource: c.Query("data-source")})
	fmt.Println(res, err)
	return c.JSON(res)
}
