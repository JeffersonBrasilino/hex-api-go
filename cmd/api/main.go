package main

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/hex-api-go/internal/user"
)

func main() {
	fmt.Println("starting api server...")
	app := fiber.New()
	user.StartWithHttp(app)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World! iii aaaa")
	})

	app.Listen(":3000")
}
