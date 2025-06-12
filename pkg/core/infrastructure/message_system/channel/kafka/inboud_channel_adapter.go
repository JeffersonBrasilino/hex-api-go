package kafka

import (
	"fmt"

	"github.com/IBM/sarama"
	"github.com/hex-api-go/pkg/core/infrastructure/message_system/container"
	"github.com/hex-api-go/pkg/core/infrastructure/message_system/message"
	"github.com/hex-api-go/pkg/core/infrastructure/message_system/channel"
)

type consumerChannelAdapterBuilder struct {
	channel.InboundChannelAdapterBuilder
	connectionReferenceName string
	topicName               string
	channelName             string
}

func NewConsumerChannelAdapterBuilder(
	connectionName string,
	topicName string,
	channelName string,
) *consumerChannelAdapterBuilder {
	return &consumerChannelAdapterBuilder{
		connectionReferenceName: connectionName,
		topicName:               topicName,
		channelName:             channelName,
	}
}

func (b *consumerChannelAdapterBuilder) Build(
	container container.Container[any, any],
) (message.ConsumerChannel, error) {
	con, err := container.Get(b.connectionReferenceName)
	if err != nil {
		return nil, fmt.Errorf(
			"[kafka-inbound-channel] connection does not exist for topic %s",
			b.topicName,
		)
	}

	consumer := con.(*connection).GetConsumer().(sarama.Consumer)
	adapter := NewInboundChannelAdapter(consumer, b.topicName)
	return adapter, nil
}

type inboundChannelAdapter struct {
	consumer  sarama.Consumer
	topicName string
}

func NewInboundChannelAdapter(
	consumer sarama.Consumer,
	topicName string,
) *inboundChannelAdapter {
	return &inboundChannelAdapter{
		consumer:  consumer,
		topicName: topicName,
	}
}

func (a *inboundChannelAdapter) Receive() (*message.Message, error) {
	return nil, nil
}

func (a *inboundChannelAdapter) Close() error {
	return nil
}

func (a *inboundChannelAdapter) Name() string {
	return a.topicName
}
