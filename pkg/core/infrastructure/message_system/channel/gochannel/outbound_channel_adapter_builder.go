package gochannel

type outboundChannelAdapterBuilder struct {
	channel *GoChannel
}

func NewOutboundChannelAdapterBuilder(channel *GoChannel) *outboundChannelAdapterBuilder {
	return &outboundChannelAdapterBuilder{
		channel,
	}
}

func (b *outboundChannelAdapterBuilder) Build() any {
	return NewOutboundChannelAdapter(b.channel)
}
