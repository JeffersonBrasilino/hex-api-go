// Package kafka provides Kafka integration for the message system.
//
// This package implements Kafka-specific channel adapters and connections for
// publishing and consuming messages through Apache Kafka. It provides outbound
// and inbound channel adapters with message translation capabilities.
//
// The InboundChannelAdapter implementation supports:
// - Kafka consumer integration for message consumption
// - Message translation between Kafka and internal formats
// - Asynchronous message processing with context support
// - Graceful shutdown and resource cleanup
package kafka

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/IBM/sarama"
	"github.com/google/uuid"
	"github.com/hex-api-go/pkg/core/infrastructure/message_system/container"
	"github.com/hex-api-go/pkg/core/infrastructure/message_system/message"
	"github.com/hex-api-go/pkg/core/infrastructure/message_system/message/channel/adapter"
)

// consumerChannelAdapterBuilder provides a builder pattern for creating
// Kafka inbound channel adapters with connection and topic configuration.
type consumerChannelAdapterBuilder struct {
	*adapter.InboundChannelAdapterBuilder[*sarama.ConsumerMessage]
	connectionReferenceName string
}

// inboundChannelAdapter implements the InboundChannelAdapter interface for Kafka,
// providing message consumption capabilities through a Kafka consumer.
type inboundChannelAdapter struct {
	consumer          sarama.Consumer
	topic             string
	messageTranslator adapter.InboundChannelMessageTranslator[*sarama.ConsumerMessage]
	channel           chan *message.Message
	ctx               context.Context
	close             context.CancelFunc
}

// NewConsumerChannelAdapterBuilder creates a new Kafka consumer channel
// adapter builder instance.
//
// Parameters:
//   - connectionReferenceName: reference name for the Kafka connection
//   - topicName: the Kafka topic to consume messages from
//   - consumerName: the consumer group name
//
// Returns:
//   - *consumerChannelAdapterBuilder: configured builder instance
func NewConsumerChannelAdapterBuilder(
	connectionReferenceName string,
	topicName string,
	consumerName string,
) *consumerChannelAdapterBuilder {
	builder := &consumerChannelAdapterBuilder{
		adapter.NewInboundChannelAdapterBuilder(
			consumerName,
			topicName,
			NewMessageTranslator(),
		),
		connectionReferenceName,
	}
	return builder
}

// NewInboundChannelAdapter creates a new Kafka inbound channel adapter instance.
//
// Parameters:
//   - consumer: the Kafka consumer for receiving messages
//   - topic: the Kafka topic name
//   - messageTranslator: translator for converting Kafka messages to internal format
//
// Returns:
//   - *inboundChannelAdapter: configured inbound channel adapter
func NewInboundChannelAdapter(
	consumer sarama.Consumer,
	topic string,
	messageTranslator adapter.InboundChannelMessageTranslator[*sarama.ConsumerMessage],
) *inboundChannelAdapter {
	ctx, cancel := context.WithCancel(context.Background())
	adp := &inboundChannelAdapter{
		consumer:          consumer,
		topic:             topic,
		messageTranslator: messageTranslator,
		channel:           make(chan *message.Message),
		ctx:               ctx,
		close:             cancel,
	}

	go adp.subscribeOnTopic()
	return adp
}

// Build constructs a Kafka inbound channel adapter from the dependency container.
//
// Parameters:
//   - container: dependency container containing required components
//
// Returns:
//   - message.InboundChannelAdapter: configured inbound channel adapter
//   - error: error if construction fails
func (c *consumerChannelAdapterBuilder) Build(
	container container.Container[any, any],
) (message.InboundChannelAdapter, error) {
	con, err := container.Get(c.connectionReferenceName)

	if err != nil {
		return nil, fmt.Errorf(
			"[kafka-inbound-channel] connection %s does not exist",
			c.connectionReferenceName,
		)
	}
	consumer := con.(*connection).GetConsumer()
	adapter := NewInboundChannelAdapter(consumer, c.ChannelName, c.MessageTranslator)
	return c.InboundChannelAdapterBuilder.BuildInboundAdapter(adapter), nil
}

// Name returns the topic name of the Kafka inbound channel adapter.
//
// Returns:
//   - string: the topic name
func (a *inboundChannelAdapter) Name() string {
	return a.topic
}

// Receive receives a message from the Kafka topic.
//
// Returns:
//   - *message.Message: the received message
//   - error: error if receiving fails or channel is closed
func (a *inboundChannelAdapter) Receive() (*message.Message, error) {
	result, hasOpen := <-a.channel
	if !hasOpen {
		return nil, errors.New("channel has not been opened")
	}
	return result, nil
}

// Close gracefully closes the Kafka inbound channel adapter and stops
// message consumption.
//
// Returns:
//   - error: error if closing fails (typically nil)
func (a *inboundChannelAdapter) Close() error {
	a.close()
	return nil
}

// subscribeOnTopic subscribes to the Kafka topic and processes incoming messages.
// This method runs in a separate goroutine and continuously polls for messages.
func (a *inboundChannelAdapter) subscribeOnTopic() {
	ticker := time.NewTicker(time.Second * 3)
	defer ticker.Stop()
	defer close(a.channel)
	var msgId int = 1
	for {
		select {
		case <-a.ctx.Done():
			return
		case <-ticker.C:
		}
		//message := a.messageTranslator.ToMessage(msg)

		msg := message.NewMessageBuilder().
			WithMessageType(message.Event).
			WithCorrelationId(uuid.New().String()).
			WithPayload(fmt.Sprintf("MESSAGE - %v", msgId)).
			Build()

		select {
		case a.channel <- msg: // Envio bem-sucedido
			//slog.Info("Message sent to internal channel.", "messageId", msgId)
			msgId++
		case <-a.ctx.Done(): // Contexto cancelado ENQUANTO esperava para enviar
			//slog.Info("Context cancelled while trying to send message. Dropping message.")
			return // Sai da goroutine
		}
	}
}
