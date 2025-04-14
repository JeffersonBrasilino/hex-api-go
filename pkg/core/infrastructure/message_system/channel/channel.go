package channel

import (
	"fmt"

	"github.com/hex-api-go/pkg/core/infrastructure/message_system/container"
	"github.com/hex-api-go/pkg/core/infrastructure/message_system/message"
	"github.com/hex-api-go/pkg/core/infrastructure/message_system/message/endpoint"
	"github.com/hex-api-go/pkg/core/infrastructure/message_system/message/router"
)

type (
	Connection interface {
		ReferenceName() string
		Connect() error
		GetProducer() any
		GetConsumer() any
	}

	BuildableChannelAdapter interface {
		Build(container.Container[any, any]) (message.MessageHandler, error)
		GetName() string
	}
)

func ConnectionReferenceName(name string) string {
	return fmt.Sprintf("channel-connection:%s", name)
}

var (
	connectionsBuilders = container.NewGenericContainer[string, Connection]()
	publisherBuilders   = container.NewGenericContainer[string, BuildableChannelAdapter]()
)

func AddChannelConnection(con Connection) {
	if connectionsBuilders.Has(con.ReferenceName()) {
		panic(
			fmt.Sprintf(
				"[channel-module] connection %s already exists",
				con.ReferenceName(),
			),
		)
	}

	connectionsBuilders.Set(con.ReferenceName(), con)
}

func buildChannelConnections(container container.Container[any, any]) {
	for _, v := range connectionsBuilders.GetAll() {
		err := v.Connect()
		if err != nil {
			panic(
				fmt.Sprintf(
					"[channel-module] %s",
					err,
				),
			)
		}
		container.Set(ConnectionReferenceName(v.ReferenceName()), v)
	}
}

func AddPublisherChannel(publisher BuildableChannelAdapter) {
	if publisherBuilders.Has(publisher.GetName()) {
		panic(
			fmt.Sprintf(
				"[channel-module] publisher %s already exists",
				publisher.GetName(),
			),
		)
	}

	publisherBuilders.Set(publisher.GetName(), publisher)
}

//TODO: ao criar o headerEnricher, o builder dele deve ser chamado aqui dentro, a fim de deixar o channel adapter limpo.
func buildPublisherChannels(container container.Container[any, any]) {
	for _, v := range publisherBuilders.GetAll() {
		adapter, err := v.Build(container)
		if err != nil {
			panic(
				fmt.Sprintf(
					"[channel-module] %s",
					err,
				),
			)
		}
		endpoint.AddGatewayBuilder(
			v.GetName(),
			endpoint.NewGatewayBuilder(v.GetName(),
				router.NewMessageRouterBuilder().
					WithRouterComponent(adapter),
			),
		)
	}
}

func Build(container container.Container[any, any]) {
	fmt.Println("build channels...")
	buildChannelConnections(container)
	buildPublisherChannels(container)
}
