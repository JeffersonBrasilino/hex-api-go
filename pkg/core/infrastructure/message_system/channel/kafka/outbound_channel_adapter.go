package kafka

import (
	"fmt"

	"github.com/IBM/sarama"
	"github.com/hex-api-go/pkg/core/infrastructure/message_system/channel"
	"github.com/hex-api-go/pkg/core/infrastructure/message_system/container"
	"github.com/hex-api-go/pkg/core/infrastructure/message_system/message"
)

type publisherChannelAdapterBuilder struct {
	connectionReferenceName string
	topicName               string
}

func NewPublisherChannelAdapterBuilder(
	connectionReferenceName string,
	topicName string,
) *publisherChannelAdapterBuilder {
	return &publisherChannelAdapterBuilder{
		connectionReferenceName: connectionReferenceName,
		topicName:               topicName,
	}
}

func (b *publisherChannelAdapterBuilder) GetName() string {
	return b.topicName
}

func (b *publisherChannelAdapterBuilder) Build(
	container container.Container[any, any],
) (message.MessageHandler, error) {
	connection, err := container.Get(channel.ConnectionReferenceName(b.connectionReferenceName))
	if err != nil {
		return nil, fmt.Errorf(
			"[kafka-outbound-channel] connection does not exist",
		)
	}

	producer := connection.(channel.Connection).GetProducer().(sarama.SyncProducer)
	adapter := NewOutboundChannelAdapter(producer, b.topicName)
	return adapter, nil
}

type outboundChannelAdapter struct {
	producer  sarama.SyncProducer
	topicName string
}

func NewOutboundChannelAdapter(
	producer sarama.SyncProducer,
	topicName string,
) *outboundChannelAdapter {
	return &outboundChannelAdapter{
		producer:  producer,
		topicName: topicName,
	}
}

func (a *outboundChannelAdapter) Handle(msg *message.Message) (*message.Message, error) {
	msg.GetHeaders().ChannelName = a.topicName
	msgTosend := FromMessage(msg)
	_, _, err := a.producer.SendMessage(msgTosend)
	return msg, err
}
