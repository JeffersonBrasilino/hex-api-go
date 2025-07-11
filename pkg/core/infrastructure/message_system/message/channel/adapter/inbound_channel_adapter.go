package adapter

import (
	"context"
	"fmt"

	messagesystem "github.com/hex-api-go/pkg/core/infrastructure/message_system"
	"github.com/hex-api-go/pkg/core/infrastructure/message_system/message"
	"github.com/hex-api-go/pkg/core/infrastructure/message_system/message/endpoint"
)

type InboundChannelMessageTranslator[T any] interface {
	ToMessage(msg T) *message.Message
}

type InboundChannelAdapterBuilder[TMessageType any] struct {
	ReferenceName         string
	ChannelName           string
	MessageTranslator     InboundChannelMessageTranslator[TMessageType]
	deadLetterChannelName string
	beforeProcessors      []message.MessageHandler
	afterProcessors       []message.MessageHandler
}

func (b *InboundChannelAdapterBuilder[TMessageType]) WithDeadLetterChannelName(
	value string,
) {
	b.deadLetterChannelName = value
}

func (b *InboundChannelAdapterBuilder[TMessageType]) WithBeforeInterceptors(
	processors ...message.MessageHandler,
) {
	b.beforeProcessors = processors
}

func (b *InboundChannelAdapterBuilder[TMessageType]) WithAfterInterceptors(
	processors ...message.MessageHandler,
) {
	b.afterProcessors = processors
}

func (b *InboundChannelAdapterBuilder[TMessageType]) BuildInboundAdapter(
	inboundAdapter message.ConsumerChannel,
) *InboundChannelAdapter {

	gatewayBuilder := endpoint.NewGatewayBuilder(b.ReferenceName, "")

	if len(b.beforeProcessors) > 0 {
		gatewayBuilder.WithBeforeInterceptors(b.beforeProcessors...)
	}

	if len(b.afterProcessors) > 0 {
		gatewayBuilder.WithAfterInterceptors(b.afterProcessors...)
	}

	messagesystem.AddGateway(gatewayBuilder)

	return NewInboundChannelAdapter(inboundAdapter)

}

type InboundChannelAdapter struct {
	inboundAdapter message.ConsumerChannel
}

func NewInboundChannelAdapter(
	adapter message.ConsumerChannel,
) *InboundChannelAdapter {
	return &InboundChannelAdapter{
		inboundAdapter: adapter,
	}
}

func (i *InboundChannelAdapter) ReceiveMessage(ctx context.Context) (*message.Message, error) {
	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("[KAFKA OUTBOUND CHANNEL] Context cancelled after processing, before sending result. ")
	default:
	}
	return i.inboundAdapter.Receive()
}

func (i *InboundChannelAdapter) Close() error {
	return i.inboundAdapter.Close()
}
