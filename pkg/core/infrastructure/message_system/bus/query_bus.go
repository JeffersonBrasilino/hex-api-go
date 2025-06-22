package bus

import (
	"context"
	"sync"

	"github.com/google/uuid"
	"github.com/hex-api-go/pkg/core/infrastructure/message_system/message"
	"github.com/hex-api-go/pkg/core/infrastructure/message_system/message/handler"
)

var createdQueryBus sync.Map

type QueryBus struct {
	*messageBus
}

func NewQueryBus(gateway message.Gateway, channelName string) *QueryBus {

	bus, ok := createdEventBus.Load(channelName)
	if ok {
		return bus.(*QueryBus)
	}
	queryBus := &QueryBus{
		messageBus: &messageBus{
			gateway,
		},
	}

	createdQueryBus.Store(channelName, bus)
	return queryBus
}

func (c *QueryBus) Send(ctx context.Context, action handler.Action) (any, error) {
	builder := c.buildMessage()
	msg := builder.WithPayload(action).
		WithRoute(action.Name()).
		Build()

	return c.sendMessage(ctx, msg)
}

func (c *QueryBus) SendRaw(ctx context.Context, route string, payload []byte, headers map[string]string) (any, error) {
	builder := c.buildMessage()
	msg := builder.WithPayload(payload).
		WithRoute(route).
		WithCustomHeader(headers).
		Build()
	return c.sendMessage(ctx, msg)
}

func (c *QueryBus) SendAsync(ctx context.Context, action handler.Action) error {
	builder := c.buildMessage()
	msg := builder.WithPayload(action).
		WithRoute(action.Name()).
		Build()
	return c.publishMessage(ctx, msg)
}

func (c *QueryBus) SendRawAsync(ctx context.Context, route string, payload any, headers map[string]string) error {
	builder := c.buildMessage()
	msg := builder.WithPayload(payload).
		WithRoute(route).
		WithCustomHeader(headers).
		Build()
	return c.publishMessage(ctx, msg)
}

func (c *QueryBus) buildMessage() *message.MessageBuilder {
	builder := message.NewMessageBuilder().
		WithMessageType(message.Query).
		WithCorrelationId(uuid.New().String())
	return builder
}
