package endpoint

import (
	"fmt"
	"log/slog"

	"github.com/hex-api-go/pkg/core/infrastructure/message_system/container"
	"github.com/hex-api-go/pkg/core/infrastructure/message_system/message"
	"github.com/hex-api-go/pkg/core/infrastructure/message_system/message/router"
)

type asyncGatewayBuilder struct {
	referenceName      string
	requestChannelName string
	beforeInterceptors []message.MessageHandler
	afterInterceptors  []message.MessageHandler
	errorChannelName   string
	replyChannelName   string
}

func NewAsyncGatewayBuilder(
	referenceName string,
	requestChannelName string,
) *asyncGatewayBuilder {
	return &asyncGatewayBuilder{
		referenceName:      referenceName,
		requestChannelName: requestChannelName,
	}
}

func (b *asyncGatewayBuilder) ReferenceName() string {
	return GatewayReferenceName(b.referenceName)
}

func (b *asyncGatewayBuilder) WithBeforeInterceptors(
	interceptors ...message.MessageHandler,
) *asyncGatewayBuilder {
	b.beforeInterceptors = append(b.beforeInterceptors, interceptors...)
	return b
}

func (b *asyncGatewayBuilder) WithAfterInterceptors(
	interceptors ...message.MessageHandler,
) *asyncGatewayBuilder {
	b.afterInterceptors = append(b.afterInterceptors, interceptors...)
	return b
}

func (b *asyncGatewayBuilder) WithErrorChannel(
	channelName string,
) *asyncGatewayBuilder {
	b.errorChannelName = channelName
	return b
}

func (b *asyncGatewayBuilder) WithReplyChannel(
	channelName string,
) *asyncGatewayBuilder {
	b.replyChannelName = channelName
	return b
}

func (b *asyncGatewayBuilder) Build(
	container container.Container[any, any],
) (message.Gateway, error) {
	requestChannel, err := container.Get(b.requestChannelName)

	if err != nil {
		panic(fmt.Sprintf("[async-gateway-builder] %s", err))
	}

	messageRouterProcessor := router.NewMessageRouterBuilder()
	messageRouterProcessor.WithRouterComponent(
		router.NewSendToChannel(
			requestChannel.(message.PublisherChannel),
		),
	)
	messageProcessor := messageRouterProcessor.Build(container)
	return NewAsyncGateway(messageProcessor, b.replyChannelName), nil
}

type AsyncGateway struct {
	messageProcessor message.MessageHandler
	replyChannelName string
}

func NewAsyncGateway(
	messageProcessor message.MessageHandler,
	replyChannelName string,
) *AsyncGateway {
	return &AsyncGateway{
		messageProcessor: messageProcessor,
		replyChannelName: replyChannelName,
	}
}

func (g *AsyncGateway) Execute(
	msg *message.Message,
) (any, error) {
	messageToProcess := message.NewMessageBuilderFromMessage(msg)
	if g.replyChannelName != "" {
		messageToProcess.WithReplyChannelName(g.replyChannelName)
	}
	fmt.Println(messageToProcess.Build().GetHeaders())
	resultMessage, err := g.messageProcessor.Handle(messageToProcess.Build())
	if err != nil {
		slog.Error("Failed to process message",
			"messageId", messageToProcess.Build().GetHeaders().MessageId,
			"reason", err.Error(),
		)
	}
	return resultMessage.GetPayload(), nil
}

func (g *AsyncGateway) IsSync() bool {
	return false
}
