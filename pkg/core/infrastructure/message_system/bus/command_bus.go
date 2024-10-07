package bus

import (
	"fmt"

	"github.com/hex-api-go/pkg/core/infrastructure/message_system/gateway"
	"github.com/hex-api-go/pkg/core/infrastructure/message_system/message"
)

type MessageSystemCommandBus struct {
	gateway *gateway.Gateway
}

func NewCommandBus() *MessageSystemCommandBus {
	return &MessageSystemCommandBus{}
}

func (b *MessageSystemCommandBus) WithChannelGateway(channelName string) *MessageSystemCommandBus {
	gatewayChannel, err := gateway.GetGateway(channelName)
	if err != nil {
		panic(
			fmt.Sprintf(
				"channel %s for command bus was not registered",
				channelName,
			),
		)
	}
	b.gateway = gatewayChannel
	return b
}

func (c *MessageSystemCommandBus) Send(
	route string,
	payload []byte,
	properties map[string]string,
) error {
	msg := c.buildMessage(route, payload, properties)
	c.gateway.Execute(msg)
	return nil
}

func (c *MessageSystemCommandBus) buildMessage(
	route string,
	payload []byte,
	properties map[string]string,
) *message.GenericMessage {
	msg := message.NewMessageBuilder()
	msg.WithPayload(payload)
	msg.WithRoute(route)
	msg.WithCustomHeader(properties)
	msg.WithMessageType(message.Command)
	msg.WithChannelName(c.gateway.Name())
	return msg.Build()
}
