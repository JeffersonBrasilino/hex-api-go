package gochannel

import "github.com/hex-api-go/pkg/core/infrastructure/message_system/message"

type outboundChannelAdapter struct {
	connection *GoChannel
}

func NewOutboundChannelAdapter(connection *GoChannel) *outboundChannelAdapter {
	return &outboundChannelAdapter{
		connection: connection,
	}
}

func (adapter *outboundChannelAdapter) Send(msg message.Message) error{
	adapter.connection.Send(msg)
	return nil
}

func (adapter *outboundChannelAdapter) Name() string{
	return adapter.connection.name
}