package http

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/hex-api-go/internal/user/infrastructure/config"
)

type httpHandlers struct {
	module *config.UserModule
}

func NewHttpHandlers(module *config.UserModule, app *fiber.App) {
	handlers := &httpHandlers{module}
	router := app.Group("/users")
	router.Post("/create", handlers.createUser)
}

func (h *httpHandlers) createUser(c *fiber.Ctx) error {
	fmt.Println(string(c.Body()))
	//h.actions.CreateUserUseCase.CreateUser("teste")
	return c.SendString("http called aaaaaaaaaaaaa")
}
