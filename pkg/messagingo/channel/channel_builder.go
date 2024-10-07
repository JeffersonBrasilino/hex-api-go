package channel

type messageChannelBuilder struct {
	inbound     InboundChannel
	outbound    OutboundChannel
	channelName string
}

/* func MessageChannelBuilder(
	name string,
	connection any,
) *messageChannelBuilder {
	return &messageChannelBuilder{
		outbound: CreateOutboundChannelAdapter(name, connection),
		inbound:  CreateInboundChannelAdapter(connection),
	}
} */

func (c *messageChannelBuilder) GetOutboundAdapter() OutboundChannel {
	return c.outbound
}

func (c *messageChannelBuilder) GetInboundAdapter() InboundChannel {
	return c.inbound
}
func (c *messageChannelBuilder) GetChannelName() string {
	return c.channelName
}
