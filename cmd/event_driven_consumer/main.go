package main

import (
	"context"
	"log/slog"
	"os/signal"
	"syscall"

	"github.com/hex-api-go/internal/user"
	"github.com/hex-api-go/pkg/core/infrastructure/messagesystem"
	kafka "github.com/hex-api-go/pkg/core/infrastructure/messagesystem/channel/kafka"
	"github.com/iyashjayesh/monigo"
)

func main() {

	monigoInstance := &monigo.Monigo{
		ServiceName:             "hex-api-go", // Mandatory field
		DashboardPort:           6061,         // Default is 8080
		DataPointsSyncFrequency: "10s",        // Default is 5 Minutes
		DataRetentionPeriod:     "1h",         // Default is 7 days. Supported values: "1h", "1d", "1w", "1m"
		TimeZone:                "Local",      // Default is Local timezone. Supported values: "Local", "UTC", "Asia/Kolkata", "America/New_York" etc. (https://en.wikipedia.org/wiki/List_of_tz_database_time_zones)
		// MaxCPUUsage:             90,         // Default is 95%
		// MaxMemoryUsage:          90,         // Default is 95%
		MaxGoRoutines: 100000, // Default is 100
	}

	go monigoInstance.Start()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	slog.Info("start message system consumer....")
	user.Bootstrap()

	messagesystem.AddChannelConnection(
		kafka.NewConnection("defaultConKafka", []string{"kafka:9092"}),
	)

	topicConsumerChannel := kafka.NewConsumerChannelAdapterBuilder(
		"defaultConKafka",
		"messagesystem.topic",
		"test_consumer",
	)
	messagesystem.AddConsumerChannel(topicConsumerChannel)

	messagesystem.Start()

	consumer, err := messagesystem.EventDrivenConsumer("test_consumer")
	if err != nil {
		panic(err)
	}

	consumer.WithAmountOfProcessors(1).
		WithMessageProcessingTimeout(50000).
		WithStopOnError(true).
		Run(ctx)

	//time.Sleep(time.Second * 10)
	//consumer.Stop()
	//stop()

	<-ctx.Done()
	messagesystem.Shutdown()
}
