package bus

import (
	"fmt"

	"github.com/hex-api-go/pkg/core/infrastructure/message_system/container"
)

type (
	BuildableBus interface {
		//GetChannelName() string provavelmente não vai usar
		GetName() string
	}
	BuildableCommandBus interface {
		BuildableBus
		// GetResponseChannelName() string não sei pq coloquei isso, tentar lembrar depois
		Build(container container.Container[any, any]) *CommandBus
	}
)

var (
	commandBusBuilders = container.NewGenericContainer[string, BuildableCommandBus]()
)

func RegisterCommandBus(busBuilder BuildableCommandBus) {
	if commandBusBuilders.Has(
		CommandBusReferenceName(busBuilder.GetName()),
	) {
		panic(
			fmt.Sprintf(
				"bus for channel %s already exists",
				CommandBusReferenceName(busBuilder.GetName()),
			),
		)
	}
	commandBusBuilders.Set(
		CommandBusReferenceName(
			busBuilder.GetName(),
		),
		busBuilder,
	)
}

func Build(container container.Container[any, any]) {
	for name, builder := range commandBusBuilders.GetAll() {
		bus := builder.Build(container)
		container.Set(name, bus)
	}
}
