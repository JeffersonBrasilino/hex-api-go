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
	referenceName      string
	requestChannelName string
}

func NewGatewayBuilder(
	referenceName string,
	requestChannelName string,
) *GatewayBuilder {
	return &GatewayBuilder{
		referenceName:      referenceName,
		requestChannelName: requestChannelName,
	}
}

func (b *GatewayBuilder) ReferenceName() string {
	return GatewayReferenceName(b.referenceName)
}

func (b *GatewayBuilder) Build(
	container container.Container[any, any],
) (message.Gateway, error) {
	requestChannel, err := container.Get(b.requestChannelName)

	if err != nil {
		panic(fmt.Sprintf("[gateway-builder] %s", err))
	}

	messageRouterProcessor := router.NewMessageRouterBuilder()
	messageRouterProcessor.WithRouterComponent(
		router.NewSendToChannel(
			requestChannel.(message.PublisherChannel),
		),
	)
	messageProcessor := messageRouterProcessor.Build(container)
	return NewGateway(messageProcessor), nil
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

	internalReplyChannel := g.makeInternalChannel()
	messageToProcess.WithReplyChannel(internalReplyChannel)

	_, err := g.messageProcessor.Handle(messageToProcess.Build())
	if err != nil {
		internalReplyChannel.Close()
		return nil, err
	}

	resultMessage := g.receive(internalReplyChannel)
	if previousReplyChannel != nil {
		previousReplyChannel.Send(resultMessage)
	}

	return resultMessage.GetPayload(), nil
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

func (g *Gateway) IsSync() bool {
	return true
}
