// Package kafka provides Kafka integration for the message system.
//
// This package implements Kafka-specific channel adapters and connections for
// publishing and consuming messages through Apache Kafka. It provides outbound
// and inbound channel adapters with message translation capabilities.
//
// The MessageTranslator implementation supports:
// - Message translation between internal and Kafka formats
// - JSON serialization and deserialization
// - Header mapping and conversion
// - Error handling for translation failures
package kafka

import (
	"encoding/json"
	"fmt"

	"github.com/IBM/sarama"
	"github.com/hex-api-go/pkg/core/infrastructure/message_system/message"
)

// MessageTranslator provides message translation capabilities between internal
// message formats and Kafka-specific formats.
type MessageTranslator struct{}

// NewMessageTranslator creates a new message translator instance.
//
// Returns:
//   - *MessageTranslator: new message translator instance
func NewMessageTranslator() *MessageTranslator {
	return &MessageTranslator{}
}

// FromMessage converts an internal message to a Kafka producer message format.
// It serializes the message headers and payload to JSON and creates appropriate
// Kafka record headers.
//
// Parameters:
//   - msg: the internal message to be converted
//
// Returns:
//   - *sarama.ProducerMessage: the Kafka producer message
func (m *MessageTranslator) FromMessage(msg *message.Message) *sarama.ProducerMessage {
	h, _ := json.Marshal(msg.GetHeaders())
	var headerMap map[string]string
	json.Unmarshal(h, &headerMap)
	saramaHeaders := []sarama.RecordHeader{}
	for k, v := range headerMap {
		saramaHeaders = append(saramaHeaders, sarama.RecordHeader{
			Key:   []byte(k),
			Value: []byte(v),
		})
	}

	payload, err := json.Marshal(msg.GetPayload())
	if err != nil {
		panic("[kafka-message-translator] cannot marshal message payload")
	}

	return &sarama.ProducerMessage{
		Topic:   msg.GetHeaders().ChannelName,
		Value:   sarama.StringEncoder(payload),
		Headers: saramaHeaders,
	}
}

// ToMessage converts a Kafka consumer message to an internal message format.
// This method is currently a placeholder and should be implemented based on
// specific requirements for message consumption.
//
// Parameters:
//   - data: the Kafka consumer message to be converted
//
// Returns:
//   - *message.Message: the internal message (placeholder implementation)
func (m *MessageTranslator) ToMessage(data *sarama.ConsumerMessage) *message.Message {
	fmt.Println("toMessage called on kafka message translator")
	return &message.Message{}
}
