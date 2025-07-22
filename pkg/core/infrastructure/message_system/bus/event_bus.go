package bus

import (
	"context"

	"github.com/google/uuid"
	"github.com/hex-api-go/pkg/core/infrastructure/message_system/message"
	"github.com/hex-api-go/pkg/core/infrastructure/message_system/message/endpoint"
	"github.com/hex-api-go/pkg/core/infrastructure/message_system/message/handler"
)

type EventBus struct {
	dispatcher *endpoint.MessageDispatcher
}

func NewEventBus(dispatcher *endpoint.MessageDispatcher) *EventBus {

	eventBus := &EventBus{
		dispatcher: dispatcher,
	}
	return eventBus
}

func (c *EventBus) Publish(ctx context.Context, action handler.Action) error {
	builder := c.buildMessage()
	msg := builder.WithPayload(action).
		WithRoute(action.Name()).
		Build()
	return c.dispatcher.PublishMessage(ctx, msg)
}

func (c *EventBus) PublishRaw(ctx context.Context, route string, payload any, headers map[string]string) error {
	builder := c.buildMessage()
	msg := builder.WithPayload(payload).
		WithRoute(route).
		WithCustomHeader(headers).
		Build()
	return c.dispatcher.PublishMessage(ctx, msg)
}

func (c *EventBus) buildMessage() *message.MessageBuilder {
	builder := message.NewMessageBuilder().
		WithMessageType(message.Event).
		WithCorrelationId(uuid.New().String())
	return builder
}
