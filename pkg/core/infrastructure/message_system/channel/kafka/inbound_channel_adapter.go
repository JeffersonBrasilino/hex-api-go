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

type consumerChannelAdapterBuilder struct {
	*adapter.InboundChannelAdapterBuilder[*sarama.ConsumerMessage]
	connectionReferenceName string
}

func NewConsumerChannelAdapterBuilder(
	connectionReferenceName string,
	topicName string,
	consumerName string,
) *consumerChannelAdapterBuilder {
	builder := &consumerChannelAdapterBuilder{
		&adapter.InboundChannelAdapterBuilder[*sarama.ConsumerMessage]{
			ReferenceName:     consumerName,
			ChannelName:       topicName,
			MessageTranslator: NewMessageTranslator(),
		},
		connectionReferenceName,
	}
	return builder
}

func (c *consumerChannelAdapterBuilder) ReferenceName() string {
	return c.InboundChannelAdapterBuilder.ReferenceName
}

func (c *consumerChannelAdapterBuilder) Build(
	container container.Container[any, any],
) (message.InboundChannelAdapter, error) {
	con, err := container.Get(c.connectionReferenceName)

	if err != nil {
		return nil, fmt.Errorf(
			"[kafka-outbound-channel] connection %s does not exist",
			c.connectionReferenceName,
		)
	}
	consumer := con.(*connection).GetConsumer()
	adapter := NewInboundChannelAdapter(consumer, c.ChannelName, c.MessageTranslator)
	return c.InboundChannelAdapterBuilder.BuildInboundAdapter(adapter), nil
}

type inboundChannelAdapter struct {
	consumer          sarama.Consumer
	topic             string
	messageTranslator adapter.InboundChannelMessageTranslator[*sarama.ConsumerMessage]
	channel           chan *message.Message
	ctx               context.Context
	close             context.CancelFunc
}

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

func (a *inboundChannelAdapter) Name() string {
	return a.topic
}

func (a *inboundChannelAdapter) Receive() (*message.Message, error) {
	result, hasOpen := <-a.channel
	if !hasOpen {
		return nil, errors.New("channel has not been opened")
	}
	return result, nil
}

func (a *inboundChannelAdapter) Close() error {
	a.close()
	return nil
}

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
