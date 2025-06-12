package kafka

import (
	"fmt"

	"github.com/IBM/sarama"
	"github.com/hex-api-go/pkg/core/infrastructure/message_system/container"
	"github.com/hex-api-go/pkg/core/infrastructure/message_system/message"
	"github.com/hex-api-go/pkg/core/infrastructure/message_system/message/channel/adapter"
)

type publisherChannelAdapterBuilder struct {
	*adapter.OutboundChannelAdapterBuilder[*sarama.ProducerMessage]
	connectionReferenceName string
}

func NewPublisherChannelAdapterBuilder(
	connectionReferenceName string,
	topicName string,
) *publisherChannelAdapterBuilder {
	builder := &publisherChannelAdapterBuilder{
		&adapter.OutboundChannelAdapterBuilder[*sarama.ProducerMessage]{},
		connectionReferenceName,
	}
	builder.WithChannelName(topicName).
		WithReferenceName(topicName).
		WithMessageTranslator(NewMessageTranslator())
	return builder
}

func (b *publisherChannelAdapterBuilder) Build(
	container container.Container[any, any],
) (message.PublisherChannel, error) {
	con, err := container.Get(b.connectionReferenceName)

	if err != nil {
		return nil, fmt.Errorf(
			"[kafka-outbound-channel] connection %s does not exist",
			b.connectionReferenceName,
		)
	}

	producer := con.(*connection).GetProducer().(sarama.SyncProducer)
	adapter := NewOutboundChannelAdapter(producer, b.ChannelName(), b.MessageTranslator())

	return b.OutboundChannelAdapterBuilder.BuildMessageHandler(adapter)
}

type outboundChannelAdapter struct {
	producer          sarama.SyncProducer
	topicName         string
	messageTranslator adapter.OutboundChannelMessageTranslator[*sarama.ProducerMessage]
}

func NewOutboundChannelAdapter(
	producer sarama.SyncProducer,
	topicName string,
	messageTranslator adapter.OutboundChannelMessageTranslator[*sarama.ProducerMessage],
) *outboundChannelAdapter {
	return &outboundChannelAdapter{
		producer:          producer,
		topicName:         topicName,
		messageTranslator: messageTranslator,
	}
}

func (a *outboundChannelAdapter) Name() string {
	return a.topicName
}

func (a *outboundChannelAdapter) Send(msg *message.Message) error {
	fmt.Println(a.topicName, " ->> send message")
	msg.GetHeaders().ChannelName = a.topicName
	msgTosend := a.messageTranslator.FromMessage(msg)
	_, _, err := a.producer.SendMessage(msgTosend)
	fmt.Println("DEU ERRO", err)
	return err
}
