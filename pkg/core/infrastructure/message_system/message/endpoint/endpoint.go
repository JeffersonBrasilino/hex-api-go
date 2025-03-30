package endpoint

import (
	"fmt"

	"github.com/hex-api-go/pkg/core/infrastructure/message_system/container"
)

type BuildableEndpoint interface {
	GetName() string
	Build(container container.Container[any, any]) error
}

var gatewayBuilders = container.NewGenericContainer[string, BuildableEndpoint]()

func AddGatewayBuilder(channelName string, builder BuildableEndpoint) {
	if gatewayBuilders.Has(channelName) {
		panic(
			fmt.Sprintf(
				"[endpoint] gateway for channel %s already exists",
				channelName,
			),
		)
	}
	gatewayBuilders.Set(channelName, builder)
}

func Build(container container.Container[any, any]) {
	fmt.Println("build endpoints...")
	for _, builder := range gatewayBuilders.GetAll() {
		err := builder.Build(container)
		if err != nil {
			panic(
				fmt.Sprintf(
					"[endpoint] %s",
					err,
				),
			)
		}
	}
}
