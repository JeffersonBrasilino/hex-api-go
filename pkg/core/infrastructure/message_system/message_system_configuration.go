package messagesystem

import (
	"fmt"

	"github.com/hex-api-go/pkg/core/infrastructure/message_system/bus"
	"github.com/hex-api-go/pkg/core/infrastructure/message_system/container"
	"github.com/hex-api-go/pkg/core/infrastructure/message_system/message/channel"
)

var (
	DefaultCommandChannelName = "message_system.default.command"
)

func registerDefaultChannels(container container.Container[any, any]) {
	chn := channel.NewDirectChannelBuilder(DefaultCommandChannelName).
		Build(container)

	container.Set(
		channel.DirectChannelReferenceName(DefaultCommandChannelName),
		chn,
	)
}

func registerDefaultBus() {
	commandBus := bus.NewCommandBusBuilder(
		DefaultCommandChannelName,
	)
	bus.RegisterCommandBus(commandBus)
}

func Build(container container.Container[any, any]) {
	fmt.Println("building default configuration...")
	registerDefaultChannels(container)
	registerDefaultBus()
}
