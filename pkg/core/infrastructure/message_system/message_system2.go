package messagesystem

import (
	"fmt"

	"github.com/hex-api-go/pkg/core/infrastructure/message_system/action"
	"github.com/hex-api-go/pkg/core/infrastructure/message_system/bus"
	"github.com/hex-api-go/pkg/core/infrastructure/message_system/channel"
	"github.com/hex-api-go/pkg/core/infrastructure/message_system/container"
	"github.com/hex-api-go/pkg/core/infrastructure/message_system/message"
	"github.com/hex-api-go/pkg/core/infrastructure/message_system/message/endpoint"
)

type (
	BuildableModule = func(container container.Container[any, any])
)

var (
	modules = []BuildableModule{
		Build,
		action.Build,
		channel.Build,
		endpoint.Build,
	}
	messageSystemContainer = container.NewGenericContainer[any, any]()
)

func buildMessageSystemModules(container container.Container[any, any]) {
	fmt.Println("load modules...")
	for _, build := range modules {
		build(container)
	}

	fmt.Println("----------------------------------------")
	fmt.Println("CONTAINER>>>", container)
	fmt.Println("----------------------------------------")
}

func Start() {
	buildMessageSystemModules(messageSystemContainer)
}

func GetCommandBus() *bus.CommandBus {
	return GetCommandBusByChannel(defaultCommandChannelName)
}

func GetQueryBus() *bus.QueryBus {
	return GetQueryBusByChannel(defaultQueryChannelName)
}

func GetCommandBusByChannel(channelName string) *bus.CommandBus {
	return bus.NewCommandBus(getBusByChannelName(channelName).(*endpoint.Gateway))
}

func GetQueryBusByChannel(channelName string) *bus.QueryBus {
	return bus.NewQueryBus(getBusByChannelName(channelName).(*endpoint.Gateway))
}

func getBusByChannelName(channelName string) any {
	found, ok := messageSystemContainer.Get(
		endpoint.GatewayReferenceName(channelName),
	)
	if ok != nil {
		panic(fmt.Sprintf("bus for channel %s not found.", channelName))
	}
	return found
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
