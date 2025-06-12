package kafka

import (
	"encoding/json"

	"github.com/IBM/sarama"
	"github.com/hex-api-go/pkg/core/infrastructure/message_system/message"
)

func ToMessage(data any) *message.Message {
	return &message.Message{}
}

type MessageTranslator struct{}

func NewMessageTranslator() *MessageTranslator {
	return &MessageTranslator{}
}

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
		panic("[kafka-message-translator] cannot marshall message payload")
	}

	return &sarama.ProducerMessage{
		Topic:   msg.GetHeaders().ChannelName,
		Value:   sarama.StringEncoder(payload),
		Headers: saramaHeaders,
	}
}
