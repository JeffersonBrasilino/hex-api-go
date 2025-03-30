package channel

import (
	"fmt"

	"github.com/hex-api-go/pkg/core/infrastructure/message_system/container"
)

type (
	Connection interface {
		ReferenceName() string
		Connect() error
		GetProducer() any
	}

	BuildableChannelAdapter interface {
		Build(container.Container[any, any]) error
		GetName() string
	}
)

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
		container.Set(v.ReferenceName(), v)
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

func buildPublisherChannels(container container.Container[any, any]) {
	for _, v := range publisherBuilders.GetAll() {
		err := v.Build(container)
		if err != nil {
			panic(
				fmt.Sprintf(
					"[channel-module] %s",
					err,
				),
			)
		}
	}
}

func Build(container container.Container[any, any]) {
	fmt.Println("build channels...")
	buildChannelConnections(container)
	buildPublisherChannels(container)
}
