package bus

import (
	"github.com/google/uuid"
	"github.com/hex-api-go/pkg/core/infrastructure/message_system/container"
	"github.com/hex-api-go/pkg/core/infrastructure/message_system/message"
	"github.com/hex-api-go/pkg/core/infrastructure/message_system/message/endpoint"
)

var createdEventBus = container.NewGenericContainer[string, *EventBus]()

type EventBus struct {
	*messageBus
}

func NewEventBus(gateway *endpoint.Gateway, channelName string) *EventBus {

	if createdBus.Has(channelName) {
		bus, _ := createdEventBus.Get(channelName)
		return bus
	}

	bus := &EventBus{
		messageBus: &messageBus{
			gateway,
		},
	}
	createdEventBus.Set(channelName, bus)
	return bus
}

func (c *EventBus) Publish(action endpoint.Action) error {
	builder := c.buildMessage()
	msg := builder.WithPayload(action).
		WithRoute(action.Name()).
		Build()
	return c.publishMessage(msg)
}

func (c *EventBus) PublishRaw(route string, payload any, headers map[string]string) error {
	builder := c.buildMessage()
	msg := builder.WithPayload(payload).
		WithRoute(route).
		WithCustomHeader(headers).
		Build()
	return c.publishMessage(msg)
}

func (c *EventBus) buildMessage() *message.MessageBuilder {
	builder := message.NewMessageBuilder().
		WithMessageType(message.Event).
		WithCorrelationId(uuid.New().String())
	return builder
}
