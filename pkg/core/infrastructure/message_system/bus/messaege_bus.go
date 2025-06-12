package bus

import (
	"fmt"

	"github.com/hex-api-go/pkg/core/infrastructure/message_system/message"
)

type messageBus struct {
	gateway message.Gateway
}

func (m *messageBus) sendMessage(msg *message.Message) (any, error) {
	if !m.gateway.IsSync() {
		return nil, fmt.Errorf("the channel is asynchronous, so it is not possible to process the message synchronously")
	}

	result, err := m.gateway.Execute(msg)
	if err != nil {
		panic(err)
	}

	resultExecution, ok := result.([]any)
	if !ok {
		panic("[message bus] got unexpected result type for message result")
	}

	if resultExecution[1] != nil {
		err, ok := resultExecution[1].(error)
		if !ok {
			panic("[message bus] unexpected result type for message")
		}
		return nil, err
	}

	return resultExecution[0], err
}

func (m *messageBus) publishMessage(msg *message.Message) error {
	if m.gateway.IsSync() {
		return fmt.Errorf("the channel is synchronously, so it is not possible to process the message asynchronous")
	}

	_, err := m.gateway.Execute(msg)
	if err != nil {
		return err
	}

	return nil
}
