package bus

import (
	"context"

	"github.com/google/uuid"
	"github.com/hex-api-go/pkg/core/infrastructure/message_system/message"
	"github.com/hex-api-go/pkg/core/infrastructure/message_system/message/endpoint"
	"github.com/hex-api-go/pkg/core/infrastructure/message_system/message/handler"
)

type CommandBus struct {
	dispatcher *endpoint.MessageDispatcher
}

func NewCommandBus(dispatcher *endpoint.MessageDispatcher) *CommandBus {
	commandBus := &CommandBus{
		dispatcher: dispatcher,
	}
	return commandBus
}

func (c *CommandBus) Send(
	ctx context.Context,
	action handler.Action,
) (any, error) {

	builder := c.buildMessage()
	msg := builder.WithPayload(action).
		WithRoute(action.Name()).
		Build()
	return c.dispatcher.SendMessage(ctx, msg)
}

func (c *CommandBus) SendRaw(
	ctx context.Context,
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

func (c *CommandBus) SendAsync(
	ctx context.Context,
	action handler.Action,
) error {
	builder := c.buildMessage()
	msg := builder.WithPayload(action).
		WithRoute(action.Name()).
		Build()
	return c.dispatcher.PublishMessage(ctx, msg)
}

func (c *CommandBus) SendRawAsync(
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

func (c *CommandBus) buildMessage() *message.MessageBuilder {
	builder := message.NewMessageBuilder().
		WithMessageType(message.Command).
		WithCorrelationId(uuid.New().String())
	return builder
}
