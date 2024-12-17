package endpoint

import (
	"fmt"

	"github.com/hex-api-go/pkg/core/infrastructure/message_system/container"
	"github.com/hex-api-go/pkg/core/infrastructure/message_system/message"
)

type (
	PollConsumer interface {
		ReceiveWithTimeout() *message.Message
	}

	ConsumerSystem interface {
		Run()
		Stop()
	}

	Endpoint interface {
		Execute(message *message.Message)
	}
)

var gatewayBuilders = container.NewGenericContainer[string, *GatewayBuilder]()

func AddGatewayBuilder(channelName string, builder *GatewayBuilder) {
	if gatewayBuilders.Has(channelName) {
		panic(
			fmt.Sprintf(
				"gateway for channel %s already exists",
				channelName,
			),
		)
	}
	gatewayBuilders.Set(channelName, builder)
}

func Build(container container.Container[any, any]) {
	for name, builder := range gatewayBuilders.GetAll() {
		gateway := builder.Build(container)
		container.Set(GatewayReferenceName(name), gateway)
	}
}

func GatewayReferenceName(referenceName string) string {
	return fmt.Sprintf("gateway:%s", referenceName)
}
