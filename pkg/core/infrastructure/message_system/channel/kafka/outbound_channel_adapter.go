// Package kafka provides Kafka integration for the message system.
//
// This package implements Kafka-specific channel adapters and connections for
// publishing and consuming messages through Apache Kafka. It provides outbound
// and inbound channel adapters with message translation capabilities.
//
// The OutboundChannelAdapter implementation supports:
// - Kafka producer integration for message publishing
// - Message translation between internal and Kafka formats
// - Context-aware message sending with timeout support
// - Error handling and connection management
package kafka

import (
	"context"
	"fmt"

	"github.com/IBM/sarama"
	"github.com/hex-api-go/pkg/core/infrastructure/message_system/container"
	"github.com/hex-api-go/pkg/core/infrastructure/message_system/message"
	"github.com/hex-api-go/pkg/core/infrastructure/message_system/message/channel/adapter"
)

// publisherChannelAdapterBuilder provides a builder pattern for creating
// Kafka outbound channel adapters with connection and topic configuration.
type publisherChannelAdapterBuilder struct {
	*adapter.OutboundChannelAdapterBuilder[*sarama.ProducerMessage]
	connectionReferenceName string
}

// outboundChannelAdapter implements the PublisherChannel interface for Kafka,
// providing message publishing capabilities through a Kafka producer.
type outboundChannelAdapter struct {
	producer          sarama.SyncProducer
	topicName         string
	messageTranslator adapter.OutboundChannelMessageTranslator[*sarama.ProducerMessage]
}

// NewPublisherChannelAdapterBuilder creates a new Kafka publisher channel
// adapter builder instance.
//
// Parameters:
//   - connectionReferenceName: reference name for the Kafka connection
//   - topicName: the Kafka topic to publish messages to
//
// Returns:
//   - *publisherChannelAdapterBuilder: configured builder instance
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

// NewOutboundChannelAdapter creates a new Kafka outbound channel adapter instance.
//
// Parameters:
//   - producer: the Kafka sync producer for sending messages
//   - topicName: the Kafka topic name
//   - messageTranslator: translator for converting internal messages to Kafka format
//
// Returns:
//   - *outboundChannelAdapter: configured outbound channel adapter
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

// Build constructs a Kafka outbound channel adapter from the dependency container.
//
// Parameters:
//   - container: dependency container containing required components
//
// Returns:
//   - message.PublisherChannel: configured publisher channel
//   - error: error if construction fails
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

// Name returns the topic name of the Kafka outbound channel adapter.
//
// Returns:
//   - string: the topic name
func (a *outboundChannelAdapter) Name() string {
	return a.topicName
}

// Send publishes a message to the Kafka topic with context support.
//
// Parameters:
//   - ctx: context for timeout/cancellation control
//   - msg: the message to be published
//
// Returns:
//   - error: error if sending fails or context is cancelled
func (a *outboundChannelAdapter) Send(ctx context.Context, msg *message.Message) error {
	msgToSend := a.messageTranslator.FromMessage(msg)
	_, _, err := a.producer.SendMessage(msgToSend)
	select {
	case <-ctx.Done():
		return fmt.Errorf("[KAFKA OUTBOUND CHANNEL] Context cancelled after processing, before sending result. ")
	default:
	}
	return err
}
