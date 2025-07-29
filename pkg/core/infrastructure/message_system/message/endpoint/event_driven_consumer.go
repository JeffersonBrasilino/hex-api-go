// Package endpoint implements the event-driven consumer pattern for message processing
// systems.
//
// This package provides a structure for consuming messages asynchronously and scalably,
// using multiple processors and integration with gateways and input channels. It
// facilitates the consumption, processing, and routing of messages in event-driven
// systems, with support for timeout, dead letter channels, and interceptors.
//
// The EventDrivenConsumer implementation supports:
// - Asynchronous message consumption with multiple concurrent processors
// - Integration with inbound channel adapters and gateways
// - Configurable processing timeouts and error handling
// - Graceful shutdown and resource cleanup
// - Dead letter channel support for failed messages
package endpoint

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/hex-api-go/pkg/core/infrastructure/message_system/container"
	"github.com/hex-api-go/pkg/core/infrastructure/message_system/message"
)

// EventDrivenConsumerBuilder is responsible for building EventDrivenConsumer instances.
// referenceName identifies the input channel to be consumed.
type EventDrivenConsumerBuilder struct {
	referenceName string
}

// EventDrivenConsumer represents an event-driven consumer.
// Manages multiple processors, processing queue, and integration with gateway and
// input channel.
type EventDrivenConsumer struct {
	referenceName                 string
	processingTimeoutMilliseconds int
	gateway                       *Gateway
	inboundChannelAdapter         message.InboundChannelAdapter
	amountOfProcessors            int
	processingQueue               chan *message.Message
	processorsWaitGroup           sync.WaitGroup
	ctx                           context.Context
	close                         context.CancelFunc
}

// NewEventDrivenConsumerBuilder creates a new EventDrivenConsumerBuilder instance.
//
// Parameters:
//   - referenceName: reference name of the input channel
//
// Returns:
//   - *EventDrivenConsumerBuilder: pointer to EventDrivenConsumerBuilder
func NewEventDrivenConsumerBuilder(referenceName string) *EventDrivenConsumerBuilder {
	return &EventDrivenConsumerBuilder{
		referenceName: referenceName,
	}
}

// NewEventDrivenConsumer creates a new EventDrivenConsumer instance.
//
// Parameters:
//   - referenceName: reference name of the input channel
//   - gateway: pointer to the associated Gateway
//   - inboundChannelAdapter: input channel adapter
//
// Returns:
//   - *EventDrivenConsumer: pointer to EventDrivenConsumer
func NewEventDrivenConsumer(
	referenceName string,
	gateway *Gateway,
	inboundChannelAdapter message.InboundChannelAdapter,
) *EventDrivenConsumer {
	ctx, cancel := context.WithCancel(context.Background())
	consumer := &EventDrivenConsumer{
		referenceName:                 referenceName,
		processingTimeoutMilliseconds: 100000,
		gateway:                       gateway,
		inboundChannelAdapter:         inboundChannelAdapter,
		amountOfProcessors:            1,
		ctx:                           ctx,
		close:                         cancel,
	}
	return consumer
}

// Build constructs an EventDrivenConsumer from the dependency container.
//
// Parameters:
//   - container: dependency container
//
// Returns:
//   - *EventDrivenConsumer: pointer to EventDrivenConsumer
//   - error: error if any occurs
func (b *EventDrivenConsumerBuilder) Build(
	container container.Container[any, any],
) (*EventDrivenConsumer, error) {

	anyChannel, err := container.Get(b.referenceName)
	if err != nil {
		panic(
			fmt.Sprintf(
				"[event-driven-consumer] consumer channel %s not found.",
				b.referenceName,
			),
		)
	}

	inboundChannel, ok := anyChannel.(message.InboundChannelAdapter)
	if !ok {
		panic(
			fmt.Sprintf(
				"[event-driven-consumer] consumer channel %s is not a consumer channel.",
				b.referenceName,
			),
		)
	}

	gatewayBuilder := NewGatewayBuilder(inboundChannel.ReferenceName(), "")
	if inboundChannel.DeadLetterChannelName() != "" {
		gatewayBuilder.WithDeadLetterChannel(inboundChannel.DeadLetterChannelName())
	}

	if len(inboundChannel.BeforeProcessors()) > 0 {
		gatewayBuilder.WithBeforeInterceptors(inboundChannel.BeforeProcessors()...)
	}

	if len(inboundChannel.AfterProcessors()) > 0 {
		gatewayBuilder.WithAfterInterceptors(inboundChannel.AfterProcessors()...)
	}

	gateway, _ := gatewayBuilder.Build(container)
	consumer := NewEventDrivenConsumer(
		b.referenceName,
		gateway,
		inboundChannel,
	)
	return consumer, nil
}

// WithMessageProcessingTimeout sets the message processing timeout in milliseconds.
//
// Parameters:
//   - milliseconds: timeout in milliseconds
//
// Returns:
//   - *EventDrivenConsumer: pointer to EventDrivenConsumer for method chaining
func (b *EventDrivenConsumer) WithMessageProcessingTimeout(
	milliseconds int,
) *EventDrivenConsumer {
	if milliseconds > 0 {
		b.processingTimeoutMilliseconds = milliseconds
	}
	return b
}

// WithAmountOfProcessors sets the number of concurrent processors.
//
// Parameters:
//   - value: number of processors
//
// Returns:
//   - *EventDrivenConsumer: pointer to EventDrivenConsumer for method chaining
func (b *EventDrivenConsumer) WithAmountOfProcessors(value int) *EventDrivenConsumer {
	if value > 1 {
		b.amountOfProcessors = value
	}
	return b
}

// Run starts processing messages received from the input channel.
//
// Parameters:
//   - ctx: context for cancellation and timeout control
//
// Returns:
//   - error: error if any occurs
func (e *EventDrivenConsumer) Run(ctx context.Context) error {
	e.processingQueue = make(chan *message.Message, e.amountOfProcessors)
	e.startProcessorsNodes(e.ctx)
	defer e.shutdown()
	for {
		select {
		case <-ctx.Done():
			if errors.Is(ctx.Err(), context.Canceled) {
				slog.Info(
					"[event-driven consumer] Context cancelled, stopping event-driven consumer.",
					"consumerName", e.referenceName,
					"error", ctx.Err(),
				)
				return nil
			}
			if errors.Is(ctx.Err(), context.DeadlineExceeded) {
				slog.Info(
					"[event-driven consumer] Deadline/Timeout exceeded, stopping event-driven consumer.",
					"consumerName", e.referenceName,
					"error", ctx.Err(),
				)
				return nil
			}
		case <-e.ctx.Done():
			slog.Info(
				"[event-driven consumer] stopping event-driven consumer",
				"consumerName", e.referenceName,
			)
			return nil
		default:
		}

		msg, err := e.inboundChannelAdapter.ReceiveMessage(e.ctx)
		if err != nil {
			return err
		}

		select {
		case e.processingQueue <- msg:
		case <-e.ctx.Done():
			slog.Info(
				"[event-driven consumer] stopping event-driven consumer",
				"consumerName", e.referenceName,
			)
			return nil
		}
	}
}

// Stop requests the consumer to stop by canceling the internal context.
func (e *EventDrivenConsumer) Stop() {
	e.close()
}

// shutdown ends processing, closes the input channel and waits for processors to finish.
func (e *EventDrivenConsumer) shutdown() {
	fmt.Println("shutdowning event-driven consumer...")
	e.inboundChannelAdapter.Close()
	close(e.processingQueue)
	e.processorsWaitGroup.Wait()
}

// startProcessorsNodes starts concurrent processors to consume messages from the queue.
//
// Parameters:
//   - ctx: context for cancellation and timeout control
func (e *EventDrivenConsumer) startProcessorsNodes(ctx context.Context) {
	for i := 0; i < e.amountOfProcessors; i++ {
		e.processorsWaitGroup.Add(1)
		go func(workerId int) {
			defer e.processorsWaitGroup.Done()
			for msg := range e.processingQueue {
				e.sendToGateway(ctx, msg, workerId)
			}
		}(i)
	}
}

// sendToGateway sends the message to the gateway for processing.
//
// Parameters:
//   - ctx: context for timeout control
//   - msg: message to be processed
//   - nodeId: processor identifier
func (e *EventDrivenConsumer) sendToGateway(
	ctx context.Context,
	msg *message.Message,
	nodeId int,
) {

	opCtx, cancel := context.WithTimeout(
		ctx,
		time.Duration(e.processingTimeoutMilliseconds)*time.Millisecond,
	)
	defer cancel()

	select {
	case <-opCtx.Done():
		return
	default:
	}

	slog.Info("[event-driven consumer] message processing started",
		"consumerName", e.referenceName,
		"nodeId", nodeId,
		"message", msg,
	)

	var err error
	time.Sleep(time.Second * 7)
	fmt.Println("processing OKOKOKOKOKOKOKOK")
	//_, err := e.gateway.Execute(opCtx, msg)
	if err != nil {
		slog.Error("[event-driven consumer] failed to process message",
			"error", err,
			"name", e.referenceName,
			"nodeId", nodeId,
			"message", msg,
		)
		return
	}

	select {
	case <-opCtx.Done():
		if errors.Is(opCtx.Err(), context.DeadlineExceeded) {
			slog.Info("[event-driven consumer] failed to process message",
				"consumerName", e.referenceName, "nodeId", nodeId,
				"error", opCtx.Err(),
			)
		}
		return
	default:
	}

	slog.Info("[event-driven consumer] message processing completed",
		"consumerName", e.referenceName,
		"nodeId", nodeId,
		"message", msg,
	)
}
