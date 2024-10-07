package gateway

import (
	"fmt"

	"github.com/hex-api-go/pkg/core/infrastructure/message_system/channel"
	"github.com/hex-api-go/pkg/core/infrastructure/message_system/message"
)

type GatewayBuilder struct {
	requestChannelName  string
	responseChannelName string
	errorChannelName    string
	interceptors        []message.MessageProcessor
}

func NewGatewayBuilder(requestChannelName string) *GatewayBuilder {
	return &GatewayBuilder{
		requestChannelName: requestChannelName,
	}
}

func (b *GatewayBuilder) WithResponseChannelName(name string) *GatewayBuilder {
	b.responseChannelName = name
	return b
}

func (b *GatewayBuilder) WithErrorChannelName(name string) *GatewayBuilder {
	b.errorChannelName = name
	return b
}

func (b *GatewayBuilder) WithInterceptors(
	interceptors ...message.MessageProcessor,
) *GatewayBuilder {
	b.interceptors = interceptors
	return b
}

func (b *GatewayBuilder) Build() *Gateway {
	fmt.Println(fmt.Printf("building gateway for %s", b.requestChannelName))
	if b.requestChannelName == b.responseChannelName {
		panic("request and response channels cannot be equal")
	}

	requestChannel, err := channel.GetOutboundChannel(b.requestChannelName)
	if err != nil {
		panic(err.Error())
	}

	var responseChannel channel.MessageOutboundChannel
	if b.responseChannelName != "" {
		responseChannel, err = channel.GetOutboundChannel(b.responseChannelName)
		if err != nil {
			panic(err.Error())
		}
	}

	buildedGateway := NewGateway(
		requestChannel,
		responseChannel,
		nil,
	)

	fmt.Println(fmt.Printf("build gateway for %s OK!", b.requestChannelName))
	return buildedGateway
}
