package bus

import (
	"encoding/json"
	"fmt"

	"github.com/hex-api-go/pkg/core/infrastructure/message_system/action"
	"github.com/hex-api-go/pkg/core/infrastructure/message_system/container"
	"github.com/hex-api-go/pkg/core/infrastructure/message_system/message"
	"github.com/hex-api-go/pkg/core/infrastructure/message_system/message/endpoint"
	"github.com/hex-api-go/pkg/core/infrastructure/message_system/message/router"
)

type CommandBusBuilder struct {
	name string
}

func NewCommandBusBuilder(name string) *CommandBusBuilder {
	return &CommandBusBuilder{
		name: name,
	}
}

func (b *CommandBusBuilder) GetName() string {
	return b.name
}

func (b *CommandBusBuilder) Build(container container.Container[any, any]) *CommandBus {

	processor := router.NewMessageRouterBuilder(container).
		WithRecipientListRouter().
		Build()

	gateway := endpoint.NewGatewayBuilder(b.name).
		WithMessageProcessor(processor).
		Build(container)

	return NewCommandBus(gateway)
}

type CommandBus struct {
	gateway *endpoint.Gateway
}

func NewCommandBus(gateway *endpoint.Gateway) *CommandBus {
	return &CommandBus{
		gateway: gateway,
	}
}

func (c *CommandBus) Send(action action.Action) (any, error) {
	msg := c.buildMessage(action)
	result, ok := c.gateway.Execute(msg).([]any)
	if !ok{
		panic("[command bus] Failed to send message")
	}

	return result[0], result[1].(error)
}

func (c *CommandBus) buildMessage(
	act action.Action,
) *message.Message {
	payload, _ := json.Marshal(act)
	msg := message.NewMessageBuilder().
		WithPayload(payload).
		WithMessageType(message.Command).
		WithRoute(action.ActionReferenceName(act.Name()))

	return msg.Build()
}

func CommandBusReferenceName(name string) string {
	return fmt.Sprintf("command-bus:%s", name)
}
