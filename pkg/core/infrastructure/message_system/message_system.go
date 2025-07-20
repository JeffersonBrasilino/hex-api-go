package messagesystem

import (
	"fmt"
	"log/slog"

	"github.com/hex-api-go/pkg/core/infrastructure/message_system/bus"
	"github.com/hex-api-go/pkg/core/infrastructure/message_system/container"
	"github.com/hex-api-go/pkg/core/infrastructure/message_system/message"
	"github.com/hex-api-go/pkg/core/infrastructure/message_system/message/endpoint"
	"github.com/hex-api-go/pkg/core/infrastructure/message_system/message/handler"
)

var (
	defaultCommandChannelName = "default.channel.command"
	defaultQueryChannelName   = "default.channel.query"
	outboundChannelBuilders   = container.NewGenericContainer[string, BuildableComponent[message.PublisherChannel]]()
	inboundChannelBuilders    = container.NewGenericContainer[string, BuildableComponent[message.InboundChannelAdapter]]()
	channelConnections        = container.NewGenericContainer[string, ChannelConnection]()
	messageSystemContainer    = container.NewGenericContainer[any, any]()
	activeEndpoints           = container.NewGenericContainer[string, any]()
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
)

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

func registerDefaultEndpoints(container container.Container[any, any]) {

	commandDispatcher, _ := endpoint.NewMessageDispatcherBuilder(
		defaultCommandChannelName,
		"",
	).Build(container)

	activeEndpoints.Set(defaultCommandChannelName, bus.NewCommandBus(commandDispatcher))

	queryDispatcher, _ := endpoint.NewMessageDispatcherBuilder(
		defaultQueryChannelName,
		"",
	).Build(container)
	activeEndpoints.Set(defaultQueryChannelName, bus.NewQueryBus(queryDispatcher))
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

func AddConsumerChannel(inboundChannel BuildableComponent[message.InboundChannelAdapter]) {
	if inboundChannelBuilders.Has(inboundChannel.ReferenceName()) {
		panic(
			fmt.Sprintf(
				"[consumer-channel] consumer for channel %s already exists",
				inboundChannel.ReferenceName(),
			),
		)
	}
	inboundChannelBuilders.Set(inboundChannel.ReferenceName(), inboundChannel)
}

func buildInboundChannels(container container.Container[any, any]) {
	for _, v := range inboundChannelBuilders.GetAll() {
		inboundChannel, err := v.Build(container)
		if err != nil {
			panic(fmt.Sprintf("[consumer-channel] %s", err))
		}
		container.Set(inboundChannel.ReferenceName(), inboundChannel)
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
	registerDefaultEndpoints(messageSystemContainer)
	buildChannelConnections(messageSystemContainer)
	buildOutboundChannels(messageSystemContainer)
	buildInboundChannels(messageSystemContainer)
}

func CommandBus() *bus.CommandBus {
	return CommandBusByChannel(defaultCommandChannelName)
}

func QueryBus() *bus.QueryBus {
	return QueryBusByChannel(defaultQueryChannelName)
}

func CommandBusByChannel(channelName string) *bus.CommandBus {
	dispatcher, err := activeEndpoints.Get(channelName)
	if err != nil {
		dispatcher, err := endpoint.NewMessageDispatcherBuilder(
			channelName,
			channelName,
		).Build(messageSystemContainer)
		if err != nil {
			panic(err)
		}

		commandBus := bus.NewCommandBus(dispatcher)
		activeEndpoints.Set(channelName, commandBus)
		return commandBus
	}

	commandDispatcher, ok := dispatcher.(*bus.CommandBus)
	if !ok {
		panic(fmt.Sprintf("channel %s is not command channel", channelName))
	}
	return commandDispatcher
}

func QueryBusByChannel(channelName string) *bus.QueryBus {
	dispatcher, err := activeEndpoints.Get(channelName)
	if err != nil {
		dispatcher, err := endpoint.NewMessageDispatcherBuilder(
			channelName,
			channelName,
		).Build(messageSystemContainer)
		if err != nil {
			panic(err)
		}

		queryBus := bus.NewQueryBus(dispatcher)
		activeEndpoints.Set(channelName, queryBus)
		return queryBus
	}

	queryDispatcher, ok := dispatcher.(*bus.QueryBus)
	if !ok {
		panic(fmt.Sprintf("channel %s is not query channel", channelName))
	}
	return queryDispatcher
}

/* func EventBusByChannel(channelName string) *bus.EventBus {
	return bus.NewEventBus(getGatewayByReference(channelName), channelName)
} */

func EventDrivenConsumer(consumerName string) (*endpoint.EventDrivenConsumer, error) {
	consumer, err := endpoint.
		NewEventDrivenConsumerBuilder(consumerName).
		Build(messageSystemContainer)

	if err != nil {
		return nil, err
	}

	return consumer, nil
}

func Shutdown() {
	slog.Info("[message-system] shutdowning...")
	for _, v := range messageSystemContainer.GetAll() {
		if inboundChannel, ok := v.(*endpoint.EventDrivenConsumer); ok {
			inboundChannel.Stop()
		}
		/* consumerChannel, ok := v.(message.ConsumerChannel)
		if ok {
			consumerChannel.Close()
			return
		}

		subscriberChannel, ok := v.(message.SubscriberChannel)
		if ok {
			subscriberChannel.Unsubscribe()
			return
		} */
	}
	slog.Info("[message-system] shutdown completed")
}
