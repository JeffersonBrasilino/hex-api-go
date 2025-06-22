package channel

import (
	"github.com/hex-api-go/pkg/core/infrastructure/message_system/message"
	"github.com/hex-api-go/pkg/core/infrastructure/message_system/message/router"
)

type InboundChannelAdapterBuilder struct {
	filters        []message.MessageHandler
	dlqChannelName string
	receiveTimeout int
	consumerName   string
}

func (i *InboundChannelAdapterBuilder) WithMessageFilter(
	filter router.FilterFunc,
) *InboundChannelAdapterBuilder {
	//i.filters = append(i.filters, router.NewMessageFilter(filter))
	return i
}

func (i *InboundChannelAdapterBuilder) WithDeadLetterChannelName(
	channelName string,
) *InboundChannelAdapterBuilder {
	i.dlqChannelName = channelName
	return i
}

func (i *InboundChannelAdapterBuilder) WithReceiveTimeout(
	value int,
) *InboundChannelAdapterBuilder {
	i.receiveTimeout = value
	return i
}

func (i *InboundChannelAdapterBuilder) WithConsumerName(
	value string,
) *InboundChannelAdapterBuilder {
	i.consumerName = value
	return i
}
