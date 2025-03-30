package kafka

import (
	"fmt"

	"github.com/IBM/sarama"
	"github.com/hex-api-go/pkg/core/infrastructure/message_system/channel"
	"github.com/hex-api-go/pkg/core/infrastructure/message_system/container"
	"github.com/hex-api-go/pkg/core/infrastructure/message_system/message"
	"github.com/hex-api-go/pkg/core/infrastructure/message_system/message/endpoint"
	"github.com/hex-api-go/pkg/core/infrastructure/message_system/message/router"
)

type publisherChannelAdapterBuilder struct {
	connectionReferenceName string
	topicName               string
	filters                 []message.MessageHandler
}

func NewPublisherChannelAdapterBuilder(
	connectionReferenceName string,
	topicName string,
) *publisherChannelAdapterBuilder {
	return &publisherChannelAdapterBuilder{
		connectionReferenceName: connectionReferenceName,
		topicName:               topicName,
		filters:                 []message.MessageHandler{},
	}
}

func (b *publisherChannelAdapterBuilder) GetName() string {
	return b.topicName
}

// move for inbound adapter...
/* func (b *publisherChannelAdapterBuilder) WithFilter(filter router.FilterFunc) *publisherChannelAdapterBuilder {
	b.filters = append(b.filters, router.NewMessageFilter(filter))
	return b
}

func (b *publisherChannelAdapterBuilder) WithDlqChannelName(channelName string) *publisherChannelAdapterBuilder {
	fmt.Println("WithDlqChannelName --->", channelName)
	return b
} */

func (b *publisherChannelAdapterBuilder) Build(
	container container.Container[any, any],
) error {
	connection, err := container.Get(ConnectionReferenceName(b.connectionReferenceName))
	if err != nil {
		return fmt.Errorf(
			"[kafka-outbound-channel] connection does not exist",
		)
	}

	producer := connection.(channel.Connection).GetProducer().(sarama.SyncProducer)
	adapter := NewOutboundChannelAdapter(producer, b.topicName)
	endpoint.AddGatewayBuilder(
		b.GetName(),
		endpoint.NewGatewayBuilder(b.GetName(),
			router.NewMessageRouterBuilder().
				WithRouterComponent(adapter),
		),
	)
	return nil
}

type outboundChannelAdapter struct {
	producer  sarama.SyncProducer
	topicName string
}

func NewOutboundChannelAdapter(
	producer sarama.SyncProducer,
	topicName string,
) *outboundChannelAdapter {
	a := &outboundChannelAdapter{
		producer:  producer,
		topicName: topicName,
	}

	return a
}

func (a *outboundChannelAdapter) Handle(msg *message.Message) (*message.Message, error) {
	_, _, err := a.producer.SendMessage(&sarama.ProducerMessage{
		Topic: a.topicName,
		Value: sarama.StringEncoder(msg.GetPayload()),
	})
	return msg, err
}
