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
)

func main() {

	log.Printf("starting api server...")
	app := fiber.New()
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	user.Bootstrap().
		WithHttpProtocol(ctx, app)

/* 	topicConsumerChannel := kafka.NewConsumerChannelAdapterBuilder(
		"defaultConKafka",
		"message_system.topic",
		"test_consumer",
	)
	messagesystem.AddConsumerChannel(topicConsumerChannel) */

	/* topicConsumerChannel2 := kafka.NewConsumerChannelAdapterBuilder(
		"defaultConKafka",
		"message_system.topic",
		"test_consumer2",
	)
	messagesystem.AddConsumerChannel(topicConsumerChannel2) */

	messagesystem.Start()

	go func() {
		log.Printf("http server listening on port 3000")
		if err := app.Listen(":3000"); err != nil {
			panic(err)
		}
	}()

	/* log.Printf("START CONSUMER......")
	a := messagesystem.PollingConsumer("test_consumer").
	WithAmountOfProcessors(1)
	go a.Run(ctx)

	go func ()  {
		time.Sleep(time.Second * 10)
		messagesystem.Shutdown()
	}() */

	/* monigoInstance := &monigo.Monigo{
		ServiceName:             "hex-api-go", // Mandatory field
		DashboardPort:           6060,         // Default is 8080
		DataPointsSyncFrequency: "10s",        // Default is 5 Minutes
		DataRetentionPeriod:     "1h",         // Default is 7 days. Supported values: "1h", "1d", "1w", "1m"
		TimeZone:                "Local",      // Default is Local timezone. Supported values: "Local", "UTC", "Asia/Kolkata", "America/New_York" etc. (https://en.wikipedia.org/wiki/List_of_tz_database_time_zones)
		// MaxCPUUsage:             90,         // Default is 95%
		// MaxMemoryUsage:          90,         // Default is 95%
		MaxGoRoutines: 100000, // Default is 100
	}

	monigoInstance.Start() */

	
	<-ctx.Done()
	log.Printf("shutdowning...")
	messagesystem.Shutdown()
	if err := app.Shutdown(); err != nil {
		log.Printf("shutting down server error: %v\n", err)
	}
	log.Printf("shutdown completed")

}
