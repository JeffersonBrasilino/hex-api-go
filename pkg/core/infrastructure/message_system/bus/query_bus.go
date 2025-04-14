package bus

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/hex-api-go/pkg/core/infrastructure/message_system/action"
	"github.com/hex-api-go/pkg/core/infrastructure/message_system/message"
	"github.com/hex-api-go/pkg/core/infrastructure/message_system/message/endpoint"
)

type QueryBus struct {
	gateway *endpoint.Gateway
}

func NewQueryBus(gateway *endpoint.Gateway) *QueryBus {
	return &QueryBus{
		gateway: gateway,
	}
}

func (c *QueryBus) Send(action action.Action) (any, error) {
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
		panic("[query bus] got unexpected result type for query result")
	}

	if resultExecution[1] != nil {
		err, ok := resultExecution[1].(error)
		if !ok {
			panic("[query bus] unexpected result type for query")
		}
		return nil, err
	}

	return resultExecution[0], err
}

func (c *QueryBus) buildMessage(
	act action.Action,
) (*message.Message, error) {
	if act.Type() != message.Query {
		return nil, fmt.Errorf("[query bus] Action %v not supported to QueryBus", act.Name())
	}
	msg := message.NewMessageBuilder().
		WithPayload(act).
		WithMessageType(message.Query).
		WithCorrelationId(uuid.New().String()).
		WithRoute(action.ActionReferenceName(act.Name()))

	return msg.Build(), nil
}
