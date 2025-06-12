package adapter

import (
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
}

func (b *OutboundChannelAdapterBuilder[TMessageType]) WithReferenceName(
	value string,
) *OutboundChannelAdapterBuilder[TMessageType] {
	b.referenceName = value
	return b
}

func (b *OutboundChannelAdapterBuilder[TMessageType]) WithChannelName(
	value string,
) *OutboundChannelAdapterBuilder[TMessageType] {
	b.channelName = value
	return b
}

func (b *OutboundChannelAdapterBuilder[TMessageType]) ReferenceName() string {
	return b.referenceName
}

func (b *OutboundChannelAdapterBuilder[TMessageType]) ChannelName() string {
	return b.channelName
}

func (b *OutboundChannelAdapterBuilder[TMessageType]) WithReplyChannelName(
	value string,
) *OutboundChannelAdapterBuilder[TMessageType] {
	b.replyChannelName = value
	return b
}

func (b *OutboundChannelAdapterBuilder[TMessageType]) ReplyChannelName(
	value string,
) string {
	return b.replyChannelName
}

func (b *OutboundChannelAdapterBuilder[TMessageType]) WithMessageTranslator(
	transator OutboundChannelMessageTranslator[TMessageType],
) *OutboundChannelAdapterBuilder[TMessageType] {
	b.messageTranslator = transator
	return b
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
		outboundHandler.Handle(msg)
	})

	gatewayBuilder := endpoint.NewAsyncGatewayBuilder(
		b.ReferenceName(),
		b.channelName,
	).
		WithReplyChannel(b.replyChannelName)

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

func (o *OutboundChannelAdapter) Handle(msg *message.Message) (*message.Message, error) {
	err := o.outboundAdapter.Send(msg)
	if err != nil {
		return nil, err
	}
	return msg, nil

}
