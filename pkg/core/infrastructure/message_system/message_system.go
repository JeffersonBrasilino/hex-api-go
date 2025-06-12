package messagesystem

//var messageSystemContainer = container.NewGenericContainer[any, any]()

/* func Start() {
	Build(messageSystemContainer)
} */

/* func GetCommandBus() *bus.CommandBus {
	return GetCommandBusByChannel(defaultCommandChannelName)
}

func GetQueryBus() *bus.QueryBus {
	return GetQueryBusByChannel(defaultQueryChannelName)
}

func GetCommandBusByChannel(channelName string) *bus.CommandBus {
	return bus.NewCommandBus(getGatewayByReference(channelName))
}

func GetQueryBusByChannel(channelName string) *bus.QueryBus {
	return bus.NewQueryBus(getGatewayByReference(channelName))
} */

/* func getGatewayByReference(referenceName string) *endpoint.Gateway {
	found, ok := messageSystemContainer.Get(endpoint.GatewayReferenceName(referenceName))
	if ok != nil {
		panic(fmt.Sprintf("bus for channel %s not found.", referenceName))
	}
	return found.(*endpoint.Gateway)
} */

/* func Shutdown() {
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
} */
