package kafka

import (
	"context"
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
		adapter.NewOutboundChannelAdapterBuilder(
			topicName,
			topicName,
			NewMessageTranslator(),
		),
		connectionReferenceName,
	}
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

	producer := con.(*connection).GetProducer()
	adapter := NewOutboundChannelAdapter(producer, b.ChannelName(), b.MessageTranslator())

	return b.OutboundChannelAdapterBuilder.BuildOutboundAdapter(adapter)
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

func (a *outboundChannelAdapter) Send(ctx context.Context, msg *message.Message) error {
	msgTosend := a.messageTranslator.FromMessage(msg)
	_, _, err := a.producer.SendMessage(msgTosend)
	select {
	case <-ctx.Done():
		return fmt.Errorf("[KAFKA OUTBOUND CHANNEL] Context cancelled after processing, before sending result. ")
	default:
	}
	return err
}
