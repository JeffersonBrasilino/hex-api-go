package channel

import "github.com/hex-api-go/pkg/messagingo/message"

type OutboundChannel interface {
	Send(message message.Message)
}

type InboundChannel interface {
	ReceiveMessage() ([]byte, error)
}

type MessageChannel interface {
	GetInboundAdapter() InboundChannel
	GetOutboundAdapter() OutboundChannel
	GetChannelName() string
}
