package kafka

import (
	"fmt"

	"github.com/IBM/sarama"
	"github.com/hex-api-go/pkg/core/infrastructure/message_system/container"
	"github.com/hex-api-go/pkg/core/infrastructure/message_system/message"
	"github.com/hex-api-go/pkg/core/infrastructure/message_system/message/router"
)

type consumerChannelAdapterBuilder struct {
	connectionReferenceName string
	topicName               string
	filters                 []message.MessageHandler
	dlqChannelName          string
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
		filters:                 []message.MessageHandler{},
		channelName:             channelName,
	}
}

func (b *consumerChannelAdapterBuilder) WithFilter(filter router.FilterFunc) *consumerChannelAdapterBuilder {
	b.filters = append(b.filters, router.NewMessageFilter(filter))
	return b
}

func (b *consumerChannelAdapterBuilder) WithDlqChannelName(channelName string) *consumerChannelAdapterBuilder {
	fmt.Println("WithDlqChannelName --->", channelName)
	b.dlqChannelName = channelName
	return b
}

func (b *consumerChannelAdapterBuilder) GetName() string {
	return b.channelName
}

func (b *consumerChannelAdapterBuilder) Build(
	container container.Container[any, any],
) error {
	/* connection, err := container.Get(channel.ConnectionReferenceName(b.connectionReferenceName))
	if err != nil {
		return fmt.Errorf(
			"[kafka-inbound-channel] connection does not exist for topic %s",
			b.topicName,
		)
	}

	consumer := connection.(channel.Connection).GetConsumer().(sarama.Consumer)
	adapter := NewInbooundChannelAdapter(consumer, b.topicName)
	endpoint.AddGatewayBuilder(
		b.GetName(),
		endpoint.NewGatewayBuilder(b.GetName(),
			router.NewMessageRouterBuilder().
				WithRouterComponent(adapter),
		),
	) */
	return nil
}

type inboundChannelAdapter struct {
	consumer  sarama.Consumer
	topicName string
}

func NewInbooundChannelAdapter(
	consumer sarama.Consumer,
	topicName string,
) *inboundChannelAdapter {
	return &inboundChannelAdapter{
		consumer:  consumer,
		topicName: topicName,
	}
}

func (a *inboundChannelAdapter) Receive() (*message.Message, error) {
	fmt.Println("MESSAGE RECEIVED")
	return nil, nil
}
func (a *inboundChannelAdapter) Close() error {
	fmt.Println("channel closed")
	return nil
}
