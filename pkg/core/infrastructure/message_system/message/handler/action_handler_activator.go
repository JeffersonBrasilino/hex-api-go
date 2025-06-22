package handler

import (
	"context"
	"fmt"

	"github.com/hex-api-go/pkg/core/infrastructure/message_system/container"
	"github.com/hex-api-go/pkg/core/infrastructure/message_system/message"
	"github.com/hex-api-go/pkg/core/infrastructure/message_system/message/channel"
)

type (
	Action interface {
		Name() string
	}

	ActionHandler[T Action, U any] interface {
		Handle(action T) (U, error)
	}
)

type ActionHandleActivatorBuilder[
	TInput Action,
	TOutput any,
] struct {
	referenceName string
	handler       ActionHandler[TInput, TOutput]
}

func NewActionHandleActivatorBuilder[
	TInput Action,
	TOutput any,
](
	referenceName string,
	handler ActionHandler[TInput, TOutput],
) *ActionHandleActivatorBuilder[TInput, TOutput] {
	return &ActionHandleActivatorBuilder[TInput, TOutput]{
		referenceName: referenceName,
		handler:       handler,
	}
}

func (b *ActionHandleActivatorBuilder[TInput, TOutput]) ReferenceName() string {
	return b.referenceName
}

func (b *ActionHandleActivatorBuilder[TInput, TOutput]) Build(container container.Container[any, any]) (message.PublisherChannel, error) {
	handlerActivator := NewActionHandlerActivator(b.handler)
	chn := channel.NewPointToPointChannel(b.referenceName)
	chn.Subscribe(func(msg *message.Message) {
		handlerActivator.Handle(msg.GetContext(), msg)
	})
	return chn, nil
}

type ActionHandleActivator[
	THandler ActionHandler[TInput, TOutput],
	TInput Action,
	TOutput any,
] struct {
	handler THandler
}

func NewActionHandlerActivator[
	THandler ActionHandler[TInput, TOutput],
	TInput Action,
	TOutput any,
](
	handler THandler,
) *ActionHandleActivator[THandler, TInput, TOutput] {
	return &ActionHandleActivator[THandler, TInput, TOutput]{
		handler: handler,
	}
}

func (c *ActionHandleActivator[THandler, TInput, TOutput]) Handle(
	ctx context.Context,
	msg *message.Message,
) (*message.Message, error) {
	action, ok := msg.GetPayload().(TInput)
	if !ok {
		return nil, fmt.Errorf("[action-handler] cannot process action: handler for action does not exists")
	}

	output, err := c.executeAction(action)

	resultMessageBuilder := message.NewMessageBuilder().
		WithChannelName(msg.GetHeaders().ReplyChannel.Name()).
		WithMessageType(message.Document)

	if err != nil {
		resultMessageBuilder.WithPayload(err)
	} else {
		resultMessageBuilder.WithPayload(output)
	}

	resultMessage := resultMessageBuilder.Build()
	if msg.GetHeaders().ReplyChannel != nil {
		msg.GetHeaders().ReplyChannel.Send(ctx, resultMessage)
	}

	return resultMessage, err
}

func (c *ActionHandleActivator[THandler, TInput, TOutput]) executeAction(
	args TInput,
) (TOutput, error) {
	result, err := c.handler.Handle(args)
	return result, err
}
