package bus

import (
	"context"

	"github.com/google/uuid"
	"github.com/hex-api-go/pkg/core/infrastructure/message_system/message"
	"github.com/hex-api-go/pkg/core/infrastructure/message_system/message/endpoint"
	"github.com/hex-api-go/pkg/core/infrastructure/message_system/message/handler"
)

type QueryBus struct {
	dispatcher *endpoint.MessageDispatcher
}

func NewQueryBus(dispatcher *endpoint.MessageDispatcher) *QueryBus {

	queryBus := &QueryBus{
		dispatcher: dispatcher,
	}
	return queryBus
}

func (c *QueryBus) Send(
	ctx context.Context,
	action handler.Action,
) (any, error) {
	builder := c.buildMessage()
	msg := builder.WithPayload(action).
		WithRoute(action.Name()).
		Build()

	return c.dispatcher.SendMessage(ctx, msg)
}

func (c *QueryBus) SendRaw(ctx context.Context,
	route string,
	payload []byte,
	headers map[string]string,
) (any, error) {
	builder := c.buildMessage()
	msg := builder.WithPayload(payload).
		WithRoute(route).
		WithCustomHeader(headers).
		Build()
	return c.dispatcher.SendMessage(ctx, msg)
}

func (c *QueryBus) SendAsync(
	ctx context.Context,
	action handler.Action,
) error {
	builder := c.buildMessage()
	msg := builder.WithPayload(action).
		WithRoute(action.Name()).
		Build()
	return c.dispatcher.PublishMessage(ctx, msg)
}

func (c *QueryBus) SendRawAsync(
	ctx context.Context,
	route string,
	payload any,
	headers map[string]string,
) error {
	builder := c.buildMessage()
	msg := builder.WithPayload(payload).
		WithRoute(route).
		WithCustomHeader(headers).
		Build()
	return c.dispatcher.PublishMessage(ctx, msg)
}

func (c *QueryBus) buildMessage() *message.MessageBuilder {
	builder := message.NewMessageBuilder().
		WithMessageType(message.Query).
		WithCorrelationId(uuid.New().String())
	return builder
}
