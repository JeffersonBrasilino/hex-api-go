package router

import (
	"context"
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
	ctx context.Context,
	msg *message.Message,
) (*message.Message, error) {
	route := r.chooseRoute(msg)
	actionChannel, err := r.messageSystemContainer.Get(route)

	if err != nil {
		return nil, fmt.Errorf("unprocessable message, channel handler for action for route %v not exists", route)
	}

	actionChannel.(message.PublisherChannel).Send(ctx, msg)

	return msg, nil
}

func (r *recipientListRouter) chooseRoute(msg *message.Message) string {
	var route string
	if msg.GetHeaders().ChannelName != "" {
		route = msg.GetHeaders().ChannelName
	}

	if msg.GetHeaders().Route != "" && route == "" {
		route = msg.GetHeaders().Route
	}
	return route
}
