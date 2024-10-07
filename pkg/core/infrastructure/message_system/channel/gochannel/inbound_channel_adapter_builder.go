package gochannel

type inboundChannelAdapterBuilder struct {
	channel *GoChannel
}

func NewInboundChannelAdapterBuilder(channel *GoChannel) *inboundChannelAdapterBuilder {
	return &inboundChannelAdapterBuilder{
		channel,
	}
}

func (b *inboundChannelAdapterBuilder) Build() any {
	return NewInboundChannelAdapter(b.channel)
}
