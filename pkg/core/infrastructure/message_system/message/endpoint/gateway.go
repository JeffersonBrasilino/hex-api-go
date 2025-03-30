package endpoint

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/hex-api-go/pkg/core/infrastructure/message_system/container"
	"github.com/hex-api-go/pkg/core/infrastructure/message_system/message"
	"github.com/hex-api-go/pkg/core/infrastructure/message_system/message/channel"
	"github.com/hex-api-go/pkg/core/infrastructure/message_system/message/router"
)

func GatewayReferenceName(referenceName string) string {
	return fmt.Sprintf("gateway:%s", referenceName)
}

type GatewayBuilder struct {
	referenceName           string
	messageProcessorBuilder *router.MessageRouterBuilder
}

func NewGatewayBuilder(referenceName string, messageProcessorBuilder *router.MessageRouterBuilder) *GatewayBuilder {
	return &GatewayBuilder{
		referenceName:           referenceName,
		messageProcessorBuilder: messageProcessorBuilder,
	}
}

func (b *GatewayBuilder) GetName() string {
	return GatewayReferenceName(b.referenceName)
}

func (b *GatewayBuilder) Build(container container.Container[any, any]) error {
	routerBuilder := b.messageProcessorBuilder.Build(container)
	buildedGateway := NewGateway(routerBuilder)
	container.Set(b.GetName(), buildedGateway)
	return nil
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
) (any, error) {

	previousReplyChannel := msg.GetHeaders().ReplyChannel
	messageToProcess := message.NewMessageBuilderFromMessage(msg)

	var internalReplyChannel *channel.PointToPointChannel
	if msg.ReplyRequired() && previousReplyChannel == nil {
		internalReplyChannel = g.makeInternalChannel()
	}

	messageToProcess.WithReplyChannel(internalReplyChannel)
	resultMessage, err := g.messageProcessor.Handle(messageToProcess.Build())
	if err != nil {
		internalReplyChannel.Close()
		return nil, err
	}

	if internalReplyChannel != nil {
		resultMessage = g.receive(internalReplyChannel)
	}

	if previousReplyChannel != nil {
		previousReplyChannel.Send(resultMessage)
	}	
	return resultMessage.GetInternalPayload(), nil
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
