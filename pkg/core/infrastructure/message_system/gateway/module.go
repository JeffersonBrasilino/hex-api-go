package gateway

import (
	"fmt"

	"github.com/hex-api-go/pkg/core/infrastructure/message_system/container"
)

var (
	gatewaysBuilders = []*GatewayBuilder{}
	gateways         = container.NewGenericContainer[string, *Gateway]()
)

func AddGatewayBuilder(builder *GatewayBuilder) {
	gatewaysBuilders = append(gatewaysBuilders, builder)
}

func Build() {
	for _, gateway := range gatewaysBuilders {
		buildedGateway := gateway.Build()
		buildedGatewayName := buildedGateway.Name()
		if gateways.Has(buildedGatewayName) {
			panic(
				fmt.Sprintf(
					"gateway %s already exists",
					buildedGatewayName,
				),
			)
		}

		gateways.Set(buildedGatewayName, buildedGateway)
	}
}

func GetGateway(gatewayChannelName string) (*Gateway, error) {
	gateway, err := gateways.Get(gatewayChannelName)
	if err != nil {
		return nil, fmt.Errorf("gateway channel %s not found", gatewayChannelName)
	}
	return gateway, nil
}
