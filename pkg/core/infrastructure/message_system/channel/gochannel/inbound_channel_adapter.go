package gochannel

type inboundChannelAdapter struct {
	connection *GoChannel
}

func NewInboundChannelAdapter(connection *GoChannel) *inboundChannelAdapter {
	return &inboundChannelAdapter{
		connection: connection,
	}
}

func (adapter *inboundChannelAdapter) Subscribe(callable func(msg any)) error {
	adapter.connection.Subscribe(callable)
	return nil
}

func (adapter *inboundChannelAdapter) Receive() (any, error) {
	return adapter.connection.Receive(), nil
}

func (adapter *inboundChannelAdapter) Unsubscribe() error {
	adapter.connection.Shutdown()
	return nil
}

func (adapter *inboundChannelAdapter) Name() string {
	return adapter.connection.name
}
