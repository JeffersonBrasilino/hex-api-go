package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
	userModule "github.com/hex-api-go/internal/user/infrastructure/config"
)

func main() {
	log.Printf("starting api server...")
	app := fiber.New()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	userModule.StartModuleWithHttpServer(ctx, app)

	go func() {
		log.Printf("http server listening on port 3000")
		if err := app.Listen(":3000"); err != nil {
			panic(err)
		}
	}()

	<-ctx.Done()
	log.Printf("shutdowning...")
	if err:= app.Shutdown();err != nil {
		log.Printf("shutting down server error: %v\n", err)
	}
	log.Printf("shutdown completed")
}
