package endpoint

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/hex-api-go/pkg/core/infrastructure/message_system/container"
	"github.com/hex-api-go/pkg/core/infrastructure/message_system/message"
	"github.com/hex-api-go/pkg/core/infrastructure/message_system/message/channel"
)

type GatewayBuilder struct {
	referenceName    string
	messageProcessor message.MessageHandler
}

func NewGatewayBuilder(referenceName string) *GatewayBuilder {
	return &GatewayBuilder{
		referenceName: referenceName,
	}
}

func (b *GatewayBuilder) WithMessageProcessor(
	processor message.MessageHandler,
) *GatewayBuilder {
	b.messageProcessor = processor
	return b
}

func (b *GatewayBuilder) GetReferenceName() string {
	return b.referenceName
}

func (b *GatewayBuilder) Build(container container.Container[any, any]) *Gateway {
	fmt.Println(fmt.Printf("building gateway for %s", b.referenceName))

	if b.messageProcessor == nil {
		panic(fmt.Sprintf("no message processor for %s", b.referenceName))
	}

	buildedGateway := NewGateway(b.messageProcessor)

	fmt.Println(fmt.Printf("build gateway for %s OK!", b.referenceName))
	return buildedGateway
}

type Gateway struct {
	messageProcessor message.MessageHandler
}

func NewGateway(
	messageProcessor message.MessageHandler,
) *Gateway {
	return &Gateway{
		messageProcessor: messageProcessor,
	}
}

func (g *Gateway) Execute(
	msg *message.Message,
) any {
	
	previousReplyChannel := msg.GetHeaders().ReplyChannel
	messageToProcess := message.NewMessageBuilderFromMessage(msg)

	var internalReplyChannel *channel.PointToPointChannel
	if msg.ReplyRequired() && previousReplyChannel == nil {
		internalReplyChannel = g.makeInternalChannel()
	}
	messageToProcess.WithReplyChannel(internalReplyChannel)
	g.messageProcessor.Handle(messageToProcess.Build())

	var resultMessage *message.Message
	if internalReplyChannel != nil {
		resultMessage = g.receive(internalReplyChannel)
	}

	if previousReplyChannel != nil {
		previousReplyChannel.Send(resultMessage)
	}
	return resultMessage.GetInternalPayload()
}

func (g *Gateway) makeInternalChannel() *channel.PointToPointChannel {
	internalChannelName := uuid.New().String()
	chn := channel.NewPointToPointChannel(internalChannelName)
	return chn
}

func (g *Gateway) receive(channel *channel.PointToPointChannel) *message.Message {
	msg, _ := channel.Receive()
	return msg.(*message.Message)
}
