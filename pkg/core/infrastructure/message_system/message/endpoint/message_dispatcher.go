package endpoint

import (
	"context"
	"fmt"

	"github.com/hex-api-go/pkg/core/infrastructure/message_system/container"
	"github.com/hex-api-go/pkg/core/infrastructure/message_system/message"
)

type messageDispatcherBuilder struct {
	referenceName      string
	requestChannelName string
}

type MessageDispatcher struct {
	gateway *Gateway
}

func NewMessageDispatcherBuilder(
	referenceName string,
	requestChannelName string,
) *messageDispatcherBuilder {
	return &messageDispatcherBuilder{
		referenceName: referenceName,
		requestChannelName: requestChannelName,
	}
}

func NewMessageDispatcher(gateway *Gateway) *MessageDispatcher {
	return &MessageDispatcher{
		gateway: gateway,
	}
}

func (b *messageDispatcherBuilder) Build(
	container container.Container[any, any],
) (*MessageDispatcher, error) {

	gateway, err := NewGatewayBuilder(
		b.referenceName,
		b.requestChannelName,
	).
		Build(container)

	if err != nil {
		return nil, fmt.Errorf("[message-dispatcher] %s", err)
	}

	dispatcher := NewMessageDispatcher(gateway)
	return dispatcher, nil
}

func (m *MessageDispatcher) SendMessage(
	ctx context.Context,
	msg *message.Message,
) (any, error) {

	result, err := m.gateway.Execute(ctx, msg)

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (m *MessageDispatcher) PublishMessage(
	ctx context.Context,
	msg *message.Message,
) error {
	_, err := m.gateway.Execute(ctx, msg)
	if err != nil {
		return err
	}

	return nil
}
