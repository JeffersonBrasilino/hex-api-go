package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	_ "net/http/pprof"

	"github.com/gofiber/fiber/v2"
	"github.com/hex-api-go/internal/user"
	messagesystem "github.com/hex-api-go/pkg/core/infrastructure/messagesystem"
	"github.com/hex-api-go/pkg/core/infrastructure/messagesystem/channel/kafka"
)

func main() {

	log.Printf("starting api server...")
	app := fiber.New()
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	user.Bootstrap().
		WithHttpProtocol(ctx, app)

	topicConsumerChannel := kafka.NewConsumerChannelAdapterBuilder(
		"defaultConKafka",
		"messagesystem.topic",
		"test_consumer",
	)
	messagesystem.AddConsumerChannel(topicConsumerChannel)

	messagesystem.Start()

	go func() {
		log.Printf("http server listening on port 3000")
		if err := app.Listen(":3000"); err != nil {
			panic(err)
		}
	}()

	messagesystem.ShowActiveEndpoints()

	/*
		a, _ := messagesystem.EventDrivenConsumer("test_consumer")

		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			log.Printf("START CONSUMER......")
			defer wg.Done()
			a.Run(ctx)
		}()

		a.WithAmountOfProcessors(1)
		wg.Wait() */
	/* 	go func() {
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
	messagesystem.Shutdown()
	if err := app.Shutdown(); err != nil {
		log.Printf("shutting down server error: %v\n", err)
	}
	log.Printf("shutdown completed")

}
