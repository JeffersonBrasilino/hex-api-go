package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	_ "net/http/pprof"

	"github.com/gofiber/fiber/v2"
	"github.com/hex-api-go/internal/user"
	messagesystem "github.com/hex-api-go/pkg/core/infrastructure/message_system"
	"github.com/iyashjayesh/monigo"
)

func main() {

	log.Printf("starting api server...")
	app := fiber.New()
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	user.Bootstrap().
		WithHttpProtocol(ctx, app)

	messagesystem.Start()

	go func() {
		log.Printf("http server listening on port 3000")
		if err := app.Listen(":3000"); err != nil {
			panic(err)
		}
	}()

	monigoInstance := &monigo.Monigo{
		ServiceName:             "hex-api-go", // Mandatory field
		DashboardPort:           6060,         // Default is 8080
		DataPointsSyncFrequency: "10s",        // Default is 5 Minutes
		DataRetentionPeriod:     "1h",         // Default is 7 days. Supported values: "1h", "1d", "1w", "1m"
		TimeZone:                "Local",      // Default is Local timezone. Supported values: "Local", "UTC", "Asia/Kolkata", "America/New_York" etc. (https://en.wikipedia.org/wiki/List_of_tz_database_time_zones)
		// MaxCPUUsage:             90,         // Default is 95%
		// MaxMemoryUsage:          90,         // Default is 95%
		MaxGoRoutines: 100000, // Default is 100
	}

	monigoInstance.Start()

	<-ctx.Done()
	log.Printf("shutdowning...")
	
	messagesystem.Shutdown()

	if err := app.Shutdown(); err != nil {
		log.Printf("shutting down server error: %v\n", err)
	}
	log.Printf("shutdown completed")

}
