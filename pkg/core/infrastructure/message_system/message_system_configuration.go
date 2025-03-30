package messagesystem

import (
	"github.com/hex-api-go/pkg/core/infrastructure/message_system/container"
	"github.com/hex-api-go/pkg/core/infrastructure/message_system/message/endpoint"
	"github.com/hex-api-go/pkg/core/infrastructure/message_system/message/router"
)

var (
	defaultCommandChannelName = "message_system.default.command"
	defaultQueryChannelName   = "message_system.default.query"
)

func registerDefaultCommandBus() {
	endpoint.AddGatewayBuilder(defaultCommandChannelName,
		endpoint.NewGatewayBuilder(defaultCommandChannelName,
			router.NewMessageRouterBuilder().
				WithRecipientListRouter(),
		),
	)
}

func registerDefaultQueryBus() {
	endpoint.AddGatewayBuilder(defaultQueryChannelName,
		endpoint.NewGatewayBuilder(defaultQueryChannelName,
			router.NewMessageRouterBuilder().
				WithRecipientListRouter(),
		),
	)
}

func Build(container container.Container[any, any]) {
	registerDefaultCommandBus()
	registerDefaultQueryBus()
}
