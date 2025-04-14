package bus

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/hex-api-go/pkg/core/infrastructure/message_system/action"
	"github.com/hex-api-go/pkg/core/infrastructure/message_system/message"
	"github.com/hex-api-go/pkg/core/infrastructure/message_system/message/endpoint"
)

type CommandBus struct {
	gateway *endpoint.Gateway
}

func NewCommandBus(gateway *endpoint.Gateway) *CommandBus {
	return &CommandBus{
		gateway: gateway,
	}
}

func (c *CommandBus) Send(action action.Action) (any, error) {
	msg, err := c.buildMessage(action)
	if err != nil {
		return nil, err
	}

	result, err := c.gateway.Execute(msg)
	if err != nil {
		panic(err)
	}

	resultExecution, ok := result.([]any)
	if !ok {
		panic("[command bus] got unexpected result type for command result")
	}

	if resultExecution[1] != nil {
		err, ok := resultExecution[1].(error)
		if !ok {
			panic("[command bus] unexpected result type for command")
		}
		return nil, err
	}

	return resultExecution[0], err
}

func (c *CommandBus) buildMessage(
	act action.Action,
) (*message.Message, error) {
	if act.Type() != message.Command {
		return nil, fmt.Errorf("[command bus] Action %v not supported to CommandBus", act.Name())
	}

	msg := message.NewMessageBuilder().
		WithPayload(act).
		WithMessageType(message.Command).
		WithCorrelationId(uuid.New().String()).
		WithRoute(action.ActionReferenceName(act.Name()))
	return msg.Build(), nil
}
