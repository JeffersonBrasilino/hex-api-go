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

type EventDrivenConsumerBuilder struct {
	referenceName string
}

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

func NewEventDrivenConsumerBuilder(
	referenceName string,
) *EventDrivenConsumerBuilder {
	return &EventDrivenConsumerBuilder{
		referenceName: referenceName,
	}
}

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

func (b *EventDrivenConsumer) WithMessageProcessingTimeout(
	milisseconds int,
) *EventDrivenConsumer {
	if milisseconds > 0 {
		b.processingTimeoutMilliseconds = milisseconds
	}
	return b
}

func (b *EventDrivenConsumer) WithAmountOfProcessors(
	value int,
) *EventDrivenConsumer {
	if value > 1 {
		b.amountOfProcessors = value
	}
	return b
}

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

func (e *EventDrivenConsumer) Stop() {
	e.close()
}

func (e *EventDrivenConsumer) shutdown() {
	e.inboundChannelAdapter.Close()
	close(e.processingQueue)
	e.processorsWaitGroup.Wait()
}

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

func (e *EventDrivenConsumer) sendToGateway(
	ctx context.Context,
	msg *message.Message,
	nodeId int,
) {

	opCtx, cancel := context.WithTimeout(ctx,
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
