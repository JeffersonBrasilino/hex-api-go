package http

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/hex-api-go/internal/user/infrastructure/config"
	"github.com/hex-api-go/pkg/core"
)

type httpHandlers struct {
	module *config.UserModule
}

func NewHttpHandlers(module *config.UserModule, app *fiber.App) {
	handlers := &httpHandlers{module}
	router := app.Group("/users")
	router.Post("/create", handlers.createUser)
	router.Get("", handlers.getUser)
}

func (h *httpHandlers) createUser(c *fiber.Ctx) error {
	res := h.module.CreateUserUseCase.CreateUser("teste")
	fmt.Println(res.GetValue())
	return c.SendString("http called aaaaaaaaaaaaa")
}

func (h *httpHandlers) getUser(c *fiber.Ctx) error {
	res := h.module.GetUserUseCase.Execute(c.Query("data-source"))

	if !res.IsSuccess() {
		if _, ok := res.GetError().(*core.NotFoundError); ok {
			return c.Status(404).JSON("not found")
		}

		if _, ok := res.GetError().(*core.DependencyError); ok {
			return c.Status(503).JSON(res.GetError())
		}
		
		return c.Status(500).JSON("internal error")
	}

	return c.JSON(res.GetValue())
}
