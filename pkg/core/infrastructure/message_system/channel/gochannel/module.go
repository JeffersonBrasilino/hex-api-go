package gochannel

/* func RegisterChannel(configuration *configuration) {
	channel := NewGoChannel(configuration.ChannelName())
	if configuration.buildPublisher {
		registerOutboundChannelAdapter(configuration, channel)
	}

	if configuration.buildConsumer {
		registerInboundChannelAdapter(channel)
	}
}

func registerOutboundChannelAdapter(configuration *configuration, buildedChannel *GoChannel) {
	adapterBuilder := NewOutboundChannelAdapterBuilder(buildedChannel)
	channel.AddChannelBuilder(adapterBuilder)

	 gatewayBuilder := gateway.NewGatewayBuilder(configuration.ChannelName())
	gateway.AddGatewayBuilder(gatewayBuilder)
}

func registerInboundChannelAdapter(buildedChannel *GoChannel) {
	adapterBuilder := NewInboundChannelAdapterBuilder(buildedChannel)
	channel.AddChannelBuilder(adapterBuilder)
}
*/