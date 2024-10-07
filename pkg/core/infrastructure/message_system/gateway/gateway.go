package gateway

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/hex-api-go/pkg/core/infrastructure/message_system/channel"
	"github.com/hex-api-go/pkg/core/infrastructure/message_system/channel/internalchannel"
	msMessage "github.com/hex-api-go/pkg/core/infrastructure/message_system/message"
)

type Gateway struct {
	requestChannel   channel.MessageOutboundChannel
	responseChannel  channel.MessageOutboundChannel
	messageProcessor msMessage.MessageProcessor
}

func NewGateway(
	requestChannel channel.MessageOutboundChannel,
	responseChannel channel.MessageOutboundChannel,
	processor msMessage.MessageProcessor,
) *Gateway {
	return &Gateway{
		requestChannel:   requestChannel,
		responseChannel:  responseChannel,
		messageProcessor: processor,
	}
}

func (g *Gateway) Execute(
	message msMessage.Message,
) any {
	messageToProcess := msMessage.NewMessageBuilderFromMessage(message)
	internalChannelName := g.registerInternalChannel()
	messageToProcess.WithReplyChannel(internalChannelName)
	g.requestChannel.Send(messageToProcess.Build())
	res := g.receive(internalChannelName)

	fmt.Println("GATEWAY > args > ", res)
	return res
}

func (g *Gateway) registerInternalChannel() string {
	internalChannelName := uuid.New().String()
	chn := internalchannel.NewPubSubChannel(internalChannelName)
	channel.AddOutboundChannel(chn)
	channel.AddInboundChannel(chn)
	return internalChannelName
}

func (g *Gateway) receive(channelName string) *msMessage.GenericMessage {
	consumerChannel, err := channel.GetInboundChannel(channelName)
	if err != nil {
		fmt.Println("error to consumer channel: ", err)
	}

	msg, _ := consumerChannel.Receive()
	consumerChannel.Unsubscribe()
	channel.RemoveInboundChannel(channelName)
	channel.RemoveOutboundChannel(channelName)
	return msg.(*msMessage.GenericMessage)
}

func (g *Gateway) Name() string {
	return g.requestChannel.Name()
}
