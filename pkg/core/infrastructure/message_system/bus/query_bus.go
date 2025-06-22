package bus

import (
	"github.com/hex-api-go/pkg/core/infrastructure/message_system/container"
	"github.com/hex-api-go/pkg/core/infrastructure/message_system/message"
)

var createdQueryBus = container.NewGenericContainer[string, *QueryBus]()

type QueryBus struct {
	*messageBus
}

func NewQueryBus(gateway message.Gateway, channelName string) *QueryBus {

	if createdQueryBus.Has(channelName) {
		bus, _ := createdQueryBus.Get(channelName)
		return bus
	}
	bus := &QueryBus{
		messageBus: &messageBus{
			gateway,
		},
	}

	createdQueryBus.Set(channelName, bus)
	return bus
}

/* func (c *QueryBus) Send(action endpoint.Action) (any, error) {
	builder := c.buildMessage()
	msg := builder.WithPayload(action).
		WithRoute(action.Name()).
		Build()

	return c.sendMessage(msg)
}

func (c *QueryBus) SendRaw(route string, payload []byte, headers map[string]string) (any, error) {
	builder := c.buildMessage()
	msg := builder.WithPayload(payload).
		WithRoute(route).
		WithCustomHeader(headers).
		Build()
	return c.sendMessage(msg)
} */

/* func (c *QueryBus) SendAsync(action endpoint.Action) error {
	builder := c.buildMessage()
	msg := builder.WithPayload(action).
		WithRoute(action.Name()).
		Build()
	return c.publishMessage(msg)
}

func (c *QueryBus) SendRawAsync(route string, payload any, headers map[string]string) error {
	builder := c.buildMessage()
	msg := builder.WithPayload(payload).
		WithRoute(route).
		WithCustomHeader(headers).
		Build()
	return c.publishMessage(msg)
}

func (c *QueryBus) buildMessage() *message.MessageBuilder {
	builder := message.NewMessageBuilder().
		WithMessageType(message.Query).
		WithCorrelationId(uuid.New().String())
	return builder
}
 */