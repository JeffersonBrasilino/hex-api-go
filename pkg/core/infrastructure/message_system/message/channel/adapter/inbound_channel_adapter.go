package adapter

import (
	"context"
	"fmt"

	"github.com/hex-api-go/pkg/core/infrastructure/message_system/message"
)

// TODO: kd as configs de DLQ? de retry? de message filter?
type InboundChannelMessageTranslator[T any] interface {
	ToMessage(msg T) *message.Message
}

type InboundChannelAdapterBuilder[TMessageType any] struct {
	ChannelName           string
	MessageTranslator     InboundChannelMessageTranslator[TMessageType]
	referenceName         string
	deadLetterChannelName string
	beforeProcessors      []message.MessageHandler
	afterProcessors       []message.MessageHandler
}

func NewInboundChannelAdapterBuilder[T any](
	referenceName string,
	channelName string,
	messageTranslator InboundChannelMessageTranslator[T],
) *InboundChannelAdapterBuilder[T] {
	return &InboundChannelAdapterBuilder[T]{
		ChannelName:       channelName,
		MessageTranslator: messageTranslator,
		referenceName:     referenceName,
		beforeProcessors:  []message.MessageHandler{},
		afterProcessors:   []message.MessageHandler{},
	}
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

func (b *InboundChannelAdapterBuilder[TMessageType]) ReferenceName() string {
	return b.ChannelName
}

func (b *InboundChannelAdapterBuilder[TMessageType]) BuildInboundAdapter(
	inboundAdapter message.ConsumerChannel,
) *InboundChannelAdapter {

	return NewInboundChannelAdapter(
		inboundAdapter,
		b.referenceName,
		b.deadLetterChannelName,
		b.beforeProcessors,
		b.afterProcessors,
	)

}

type InboundChannelAdapter struct {
	inboundAdapter        message.ConsumerChannel
	referenceName         string
	deadLetterChannelName string
	beforeProcessors      []message.MessageHandler
	afterProcessors       []message.MessageHandler
}

func NewInboundChannelAdapter(
	adapter message.ConsumerChannel,
	referenceName string,
	deadLetterChannelName string,
	beforeProcessors []message.MessageHandler,
	afterProcessors []message.MessageHandler,
) *InboundChannelAdapter {
	return &InboundChannelAdapter{
		inboundAdapter:        adapter,
		referenceName:         referenceName,
		deadLetterChannelName: deadLetterChannelName,
		beforeProcessors:      beforeProcessors,
		afterProcessors:       afterProcessors,
	}
}

func (i *InboundChannelAdapter) ReferenceName() string {
	return i.referenceName
}

func (i *InboundChannelAdapter) DeadLetterChannelName() string {
	return i.deadLetterChannelName
}

func (i *InboundChannelAdapter) AfterProcessors() []message.MessageHandler {
	return i.afterProcessors
}

func (i *InboundChannelAdapter) BeforeProcessors() []message.MessageHandler {
	return i.beforeProcessors
}

func (i *InboundChannelAdapter) ReceiveMessage(ctx context.Context) (*message.Message, error) {
	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("[inbound-channel] Context cancelled after processing, before sending result. ")
	default:
	}
	return i.inboundAdapter.Receive()
}

func (i *InboundChannelAdapter) Close() error {
	return i.inboundAdapter.Close()
}
