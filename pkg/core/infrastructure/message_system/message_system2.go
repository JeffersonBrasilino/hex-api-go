package messagesystem

import (
	"fmt"

	"github.com/hex-api-go/pkg/core/infrastructure/message_system/action"
	"github.com/hex-api-go/pkg/core/infrastructure/message_system/bus"
	"github.com/hex-api-go/pkg/core/infrastructure/message_system/container"
)

type (
	BuildableModule = func(container container.Container[any, any])
)

var (
	modules = []BuildableModule{
		Build,
		bus.Build,
		action.Build,
	}
	messageSystemContainer = container.NewGenericContainer[any, any]()
)

func buildMessageSystemModules(container container.Container[any, any]) {
	fmt.Println("load modules...")
	for _, build := range modules {
		build(container)
	}
}

func Start() {
	buildMessageSystemModules(messageSystemContainer)
}

func GetCommandBus() *bus.CommandBus {
	found, ok := messageSystemContainer.Get(bus.CommandBusReferenceName(DefaultCommandChannelName))
	if ok != nil {
		panic(fmt.Sprintf("commandbus not found: %v", DefaultCommandChannelName))
	}
	return found.(*bus.CommandBus)
}