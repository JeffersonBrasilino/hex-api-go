package bus

import (
	"context"

	"github.com/hex-api-go/pkg/core/infrastructure/message_system/message"
)

type messageBus struct {
	gateway message.Gateway
}

func (m *messageBus) sendMessage(ctx context.Context, msg *message.Message) (any, error) {

	result, err := m.gateway.Execute(ctx, msg)

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (m *messageBus) publishMessage(ctx context.Context, msg *message.Message) error {
	_, err := m.gateway.Execute(ctx, msg)
	if err != nil {
		return err
	}

	return nil
}
