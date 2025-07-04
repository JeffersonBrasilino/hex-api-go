package messagesystem

import (
	"fmt"

	"github.com/hex-api-go/pkg/core/infrastructure/message_system/bus"
	"github.com/hex-api-go/pkg/core/infrastructure/message_system/container"
	"github.com/hex-api-go/pkg/core/infrastructure/message_system/message"
	"github.com/hex-api-go/pkg/core/infrastructure/message_system/message/endpoint"
	"github.com/hex-api-go/pkg/core/infrastructure/message_system/message/handler"
)

var (
	defaultCommandChannelName = "default.channel.command"
	defaultQueryChannelName   = "default.channel.query"
	gatewayBuilders           = container.NewGenericContainer[string, BuildableComponent[message.Gateway]]()
	outboundChannelBuilders   = container.NewGenericContainer[string, BuildableComponent[message.PublisherChannel]]()
	channelConnections        = container.NewGenericContainer[string, ChannelConnection]()
	messageSystemContainer    = container.NewGenericContainer[any, any]()
)

type (
	ChannelConnection interface {
		ReferenceName() string
		Connect() error
		Disconnect() error
	}
	BuildableComponent[T any] interface {
		Build(container container.Container[any, any]) (T, error)
		ReferenceName() string
	}
	RegistrableComponent interface {
		Register(container container.Container[any, any])
	}
)

func AddGateway(builder BuildableComponent[message.Gateway]) {
	if gatewayBuilders.Has(builder.ReferenceName()) {
		panic(
			fmt.Sprintf(
				"[endpoint] gateway for channel %s already exists",
				builder.ReferenceName(),
			),
		)
	}
	gatewayBuilders.Set(builder.ReferenceName(), builder)
}

func buildGateways(container container.Container[any, any]) {
	for _, builder := range gatewayBuilders.GetAll() {
		gateway, err := builder.Build(container)
		if err != nil {
			panic(
				fmt.Sprintf(
					"[endpoint] %s",
					err,
				),
			)
		}
		container.Set(builder.ReferenceName(), gateway)
	}
}

func AddPublisherChannel(publisher BuildableComponent[message.PublisherChannel]) {
	if outboundChannelBuilders.Has(publisher.ReferenceName()) {
		panic(
			fmt.Sprintf(
				"[publisher-channel] channel %s already exists",
				publisher.ReferenceName(),
			),
		)
	}

	outboundChannelBuilders.Set(publisher.ReferenceName(), publisher)
}

func buildOutboundChannels(container container.Container[any, any]) {
	for _, v := range outboundChannelBuilders.GetAll() {
		outboundChannel, err := v.Build(container)
		if err != nil {
			panic(
				fmt.Sprintf(
					"[publisher-channel] %s",
					err,
				),
			)
		}

		container.Set(v.ReferenceName(), outboundChannel)

	}
}

func registerDefaultCommandBus() {
	AddGateway(endpoint.NewGatewayBuilder(
		defaultCommandChannelName,
		"",
	))

}

func registerDefaultQueryBus() {
	AddGateway(endpoint.NewGatewayBuilder(
		defaultQueryChannelName,
		"",
	))
}

func AddChannelConnection(con ChannelConnection) {
	if channelConnections.Has(con.ReferenceName()) {
		panic(
			fmt.Sprintf(
				"[channel-module] connection %s already exists",
				con.ReferenceName(),
			),
		)
	}

	channelConnections.Set(con.ReferenceName(), con)
}

func buildChannelConnections(container container.Container[any, any]) {
	for _, v := range channelConnections.GetAll() {
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

func AddActionHandler[
	T handler.Action,
	U any,
](handlerAction handler.ActionHandler[T, U]) {
	action := *new(T)
	if outboundChannelBuilders.Has(action.Name()) {
		panic(
			fmt.Sprintf(
				"hander for %s already exists",
				action.Name(),
			),
		)
	}

	outboundChannelBuilders.Set(
		action.Name(),
		handler.NewActionHandleActivatorBuilder(
			action.Name(),
			handlerAction,
		),
	)
}

func Start() {
	registerDefaultCommandBus()
	registerDefaultQueryBus()
	buildChannelConnections(messageSystemContainer)
	buildOutboundChannels(messageSystemContainer)
	buildGateways(messageSystemContainer)
}

func CommandBus() *bus.CommandBus {
	return CommandBusByChannel(defaultCommandChannelName)
}

func QueryBus() *bus.QueryBus {
	return QueryBusByChannel(defaultQueryChannelName)
}

func CommandBusByChannel(channelName string) *bus.CommandBus {
	return bus.NewCommandBus(getGatewayByReference(channelName), channelName)
}

func QueryBusByChannel(channelName string) *bus.QueryBus {
	return bus.NewQueryBus(getGatewayByReference(channelName), channelName)
}

func EventBusByChannel(channelName string) *bus.EventBus {
	return bus.NewEventBus(getGatewayByReference(channelName), channelName)
}

func getGatewayByReference(referenceName string) message.Gateway {
	found, ok := messageSystemContainer.Get(endpoint.GatewayReferenceName(referenceName))
	if ok != nil {
		panic(fmt.Sprintf("bus for channel %s not found.", referenceName))
	}
	gtw, _ := found.(message.Gateway)
	return gtw
}

func Shutdown() {
	for _, v := range messageSystemContainer.GetAll() {
		consumerChannel, ok := v.(message.ConsumerChannel)
		if ok {
			consumerChannel.Close()
		}

		subscriberChannel, ok := v.(message.SubscriberChannel)
		if ok {
			subscriberChannel.Unsubscribe()
		}
	}
}
