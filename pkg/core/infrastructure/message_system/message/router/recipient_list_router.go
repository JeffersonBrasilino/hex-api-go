package router

import (
	"fmt"

	"github.com/hex-api-go/pkg/core/infrastructure/message_system/container"
	"github.com/hex-api-go/pkg/core/infrastructure/message_system/message"
)

type recipientListRouter struct {
	messageSystemContainer container.Container[any, any]
}

func NewRecipientListRouter(
	messageSystemContainer container.Container[any, any],
) *recipientListRouter {
	return &recipientListRouter{messageSystemContainer: messageSystemContainer}
}

func (r *recipientListRouter) Handle(
	msg *message.Message,
) (*message.Message, error) {
	route := msg.GetHeaders().Route
	if route == "" {
		return nil, fmt.Errorf("unprocessable message, missing route param from header message")
	}

	actionChannel, err := r.messageSystemContainer.Get(route)
	if err != nil {
		return nil, fmt.Errorf("unprocessable message, channel handler for action for route %v not exists", route)
	}

	actionChannel.(message.PublisherChannel).Send(msg)

	return msg, nil
}
