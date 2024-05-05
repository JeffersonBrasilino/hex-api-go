package user

import (
	"github.com/gofiber/fiber/v2"
	"github.com/hex-api-go/internal/user/infrastructure/config"
	"github.com/hex-api-go/internal/user/infrastructure/http"
)

func StartWithHttp(fiberApp *fiber.App) {
	module := config.Bootstrap()
	http.NewHttpHandlers(module, fiberApp)
}

func StartWithEda() {}
