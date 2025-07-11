package endpoint

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/hex-api-go/pkg/core/infrastructure/message_system/container"
	"github.com/hex-api-go/pkg/core/infrastructure/message_system/message"
)

var startedConsumers sync.Map

type PollingConsumerBuilder struct {
	referenceName string
}

func NewPolllingConsumerBuilder(
	referenceName string,
) *PollingConsumerBuilder {
	return &PollingConsumerBuilder{
		referenceName: referenceName,
	}
}

func (b *PollingConsumerBuilder) Build(container container.Container[any, any]) (*PollingConsumer, error) {

	_, hasExists := startedConsumers.Load(b.referenceName)
	if hasExists {
		return nil, fmt.Errorf("consumer %s already started", b.referenceName)
	}

	channel, ok := container.Get(b.referenceName)
	if ok != nil {
		panic(fmt.Sprintf("consumer channel %s not found.", b.referenceName))
	}

	inboundChannel, instance := channel.(message.InboundChannelAdapter)
	if !instance {
		panic(fmt.Sprintf("consumer channel %s is not a consumer channel.", b.referenceName))
	}

	gateway, err := container.Get(GatewayReferenceName(b.referenceName))
	if err != nil {
		return nil, fmt.Errorf(
			"[polling-consumer] gateway %s does not exist",
			b.referenceName,
		)
	}

	startedConsumers.Store(b.referenceName, true)

	return NewPollingConsumer(
		gateway.(*Gateway),
		inboundChannel,
		b.referenceName,
	), nil
}

type PollingConsumer struct {
	referenceName                 string
	pollIntervalMilliseconds      int
	processingDelayMilliseconds   int
	processingTimeoutMilliseconds int
	stopOnError                   bool
	hasRunning                    bool
	gateway                       *Gateway
	inboundChannelAdapter         message.InboundChannelAdapter
}

func NewPollingConsumer(
	gateway *Gateway,
	inboundChannelAdapter message.InboundChannelAdapter,
	referenceName string,
) *PollingConsumer {
	return &PollingConsumer{
		pollIntervalMilliseconds:      1000,
		processingDelayMilliseconds:   0,
		processingTimeoutMilliseconds: 100000,
		stopOnError:                   true,
		gateway:                       gateway,
		inboundChannelAdapter:         inboundChannelAdapter,
		referenceName:                 referenceName,
	}
}

func (b *PollingConsumer) WithPollIntervalMilliseconds(
	value int,
) *PollingConsumer {
	b.pollIntervalMilliseconds = value
	return b
}

func (b *PollingConsumer) WithProcessingDelayMilliseconds(
	value int,
) *PollingConsumer {
	b.processingDelayMilliseconds = value
	return b
}

func (b *PollingConsumer) WithStopOnError(
	value bool,
) *PollingConsumer {
	b.stopOnError = value
	return b
}

func (b *PollingConsumer) WithProcessingTimeoutMilliseconds(
	value int,
) *PollingConsumer {
	b.processingTimeoutMilliseconds = value
	return b
}

func (c *PollingConsumer) Run(ctx context.Context) error {
	slog.Info("Starting polling consumer", "consumerName", c.referenceName)
	c.hasRunning = true

	ticker := time.NewTicker(time.Millisecond * time.Duration(c.pollIntervalMilliseconds))
	defer ticker.Stop()

	for c.hasRunning {
		select {
		case <-ctx.Done():
			c.Stop()
			return ctx.Err()
		case <-ticker.C:
			msg, err := c.inboundChannelAdapter.ReceiveMessage(ctx)
			if err != nil {
				slog.Error("Error receiving message", "error", err, "name", c.referenceName)
				if c.stopOnError {
					c.Stop()
					return err
				}
				continue
			}
			if msg == nil {
				slog.Info("no message received", "consumerName", c.referenceName)
				continue
			}

			if c.processingDelayMilliseconds > 0 {
				time.Sleep(time.Millisecond * time.Duration(c.processingDelayMilliseconds))
			}

			go c.sendToGateway(ctx, msg)
		}
	}

	return nil
}

func (c *PollingConsumer) sendToGateway(ctx context.Context, msg *message.Message) {
	fmt.Println("sendToGateway", msg)
	opCtx, cancel := context.WithTimeout(ctx,
		time.Duration(c.processingTimeoutMilliseconds)*time.Millisecond,
	)
	defer cancel()
	_, err := c.gateway.Execute(opCtx, msg)
	if err != nil {
		slog.Error("failed to process message",
			"error", err,
			"name", c.referenceName,
			"messageId", msg.GetHeaders().MessageId,
		)
		return
	}
	slog.Debug("message processed", "name", c.referenceName)
}

func (c *PollingConsumer) Stop() {
	c.hasRunning = false
	startedConsumers.Delete(c.referenceName)
}
