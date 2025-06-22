package adapter

import (
	"context"

	messagesystem "github.com/hex-api-go/pkg/core/infrastructure/message_system"
	"github.com/hex-api-go/pkg/core/infrastructure/message_system/message"
	"github.com/hex-api-go/pkg/core/infrastructure/message_system/message/channel"
	"github.com/hex-api-go/pkg/core/infrastructure/message_system/message/endpoint"
)

type OutboundChannelMessageTranslator[T any] interface {
	FromMessage(msg *message.Message) T
}

type OutboundChannelAdapterBuilder[TMessageType any] struct {
	referenceName     string
	channelName       string
	replyChannelName  string
	messageTranslator OutboundChannelMessageTranslator[TMessageType]
	beforeProcessors  []message.MessageHandler
	afterProcessors   []message.MessageHandler
}

func (b *OutboundChannelAdapterBuilder[TMessageType]) WithReferenceName(
	value string,
) *OutboundChannelAdapterBuilder[TMessageType] {
	b.referenceName = value
	return b
}

func (b *OutboundChannelAdapterBuilder[TMessageType]) WithChannelName(
	value string,
) {
	b.channelName = value
}

func (b *OutboundChannelAdapterBuilder[TMessageType]) WithMessageTranslator(
	transator OutboundChannelMessageTranslator[TMessageType],
) {
	b.messageTranslator = transator
}

func (b *OutboundChannelAdapterBuilder[TMessageType]) WithReplyChannelName(
	value string,
) {
	b.replyChannelName = value
}

func (b *OutboundChannelAdapterBuilder[TMessageType]) WithBeforeInterceptors(
	processors ...message.MessageHandler,
) {
	b.beforeProcessors = processors
}

func (b *OutboundChannelAdapterBuilder[TMessageType]) WithAfterInterceptors(
	processors ...message.MessageHandler,
) {
	b.afterProcessors = processors
}

func (b *OutboundChannelAdapterBuilder[TMessageType]) ReferenceName() string {
	return b.referenceName
}

func (b *OutboundChannelAdapterBuilder[TMessageType]) ChannelName() string {
	return b.channelName
}

func (b *OutboundChannelAdapterBuilder[TMessageType]) ReplyChannelName(
	value string,
) string {
	return b.replyChannelName
}

func (
	b *OutboundChannelAdapterBuilder[TMessageType],
) MessageTranslator() OutboundChannelMessageTranslator[TMessageType] {
	return b.messageTranslator
}

func (b *OutboundChannelAdapterBuilder[TMessageType]) BuildMessageHandler(
	outboundAdapter message.PublisherChannel,
) (*channel.PointToPointChannel, error) {

	outboundHandler := NewOutboundChannelAdapter(outboundAdapter)

	chn := channel.NewPointToPointChannel(b.referenceName)
	chn.Subscribe(func(msg *message.Message) {
		outboundHandler.Handle(msg.GetContext(), msg)
	})

	gatewayBuilder := endpoint.NewGatewayBuilder(
		b.ReferenceName(),
		b.channelName,
	).
		WithReplyChannel(b.replyChannelName)

	if len(b.beforeProcessors) > 0 {
		gatewayBuilder.WithBeforeInterceptors(b.beforeProcessors...)
	}

	if len(b.afterProcessors) > 0 {
		gatewayBuilder.WithAfterInterceptors(b.afterProcessors...)
	}

	messagesystem.AddGateway(gatewayBuilder)

	return chn, nil
}

type OutboundChannelAdapter struct {
	outboundAdapter message.PublisherChannel
}

func NewOutboundChannelAdapter(
	adapter message.PublisherChannel,
) *OutboundChannelAdapter {
	return &OutboundChannelAdapter{
		outboundAdapter: adapter,
	}
}

func (o *OutboundChannelAdapter) Handle(ctx context.Context, msg *message.Message) (*message.Message, error) {
	err := o.outboundAdapter.Send(ctx, msg)
	if msg.GetHeaders().ReplyChannel != nil {
		o.publishOnInternalChannel(ctx, msg, err)
	}

	if err != nil {
		return nil, err
	}

	return msg, nil
}

func (o *OutboundChannelAdapter) publishOnInternalChannel(ctx context.Context, msg *message.Message, response any) {
	payloadMessage := msg.GetPayload()
	if response != nil {
		payloadMessage = response
	}
	resultMessage := message.NewMessageBuilderFromMessage(msg).
		WithMessageType(message.Document).
		WithPayload(payloadMessage).
		Build()
	msg.GetHeaders().ReplyChannel.Send(ctx, resultMessage)
}
