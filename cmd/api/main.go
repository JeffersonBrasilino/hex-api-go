package main

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	userModule "github.com/hex-api-go/internal/user/infrastructure/config"
)

func main() {
	fmt.Println("starting api server...")
	app := fiber.New()

	userModule.Bootstrap()
	userModule.WithHttpHandlers(app)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World! iii aaaa")
	})

	app.Listen(":3000")
}
