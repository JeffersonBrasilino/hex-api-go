package endpoint

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"github.com/hex-api-go/pkg/core/infrastructure/message_system/container"
	"github.com/hex-api-go/pkg/core/infrastructure/message_system/message"
	"github.com/hex-api-go/pkg/core/infrastructure/message_system/message/channel"
	"github.com/hex-api-go/pkg/core/infrastructure/message_system/message/handler"
	"github.com/hex-api-go/pkg/core/infrastructure/message_system/message/router"
)
func GatewayReferenceName(referenceName string) string {
	return fmt.Sprintf("gateway:%s", referenceName)
}

type gatewayBuilder struct {
	referenceName      string
	requestChannelName string
	beforeInterceptors []message.MessageHandler
	afterInterceptors  []message.MessageHandler
	deadLetterChannel  string
	replyChannelName   string
}

func NewGatewayBuilder(
	referenceName string,
	requestChannelName string,
) *gatewayBuilder {
	return &gatewayBuilder{
		referenceName:      referenceName,
		requestChannelName: requestChannelName,
	}
}

func (b *gatewayBuilder) ReferenceName() string {
	return GatewayReferenceName(b.referenceName)
}

func (b *gatewayBuilder) WithBeforeInterceptors(
	interceptors ...message.MessageHandler,
) *gatewayBuilder {
	b.beforeInterceptors = append(b.beforeInterceptors, interceptors...)
	return b
}

func (b *gatewayBuilder) WithAfterInterceptors(
	interceptors ...message.MessageHandler,
) *gatewayBuilder {
	b.afterInterceptors = append(b.afterInterceptors, interceptors...)
	return b
}

func (b *gatewayBuilder) WithDeadLetterChannel(
	channelName string,
) *gatewayBuilder {
	b.deadLetterChannel = channelName
	return b
}

func (b *gatewayBuilder) WithReplyChannel(
	channelName string,
) *gatewayBuilder {
	b.replyChannelName = channelName
	return b
}

func (b *gatewayBuilder) Build(
	container container.Container[any, any],
) (*Gateway, error) {

	rt := router.NewRouter()
	if b.beforeInterceptors != nil {
		for _, beforeInterceptors := range b.beforeInterceptors {
			rt.AddHandler(handler.NewContextHandler(beforeInterceptors))
		}
	}

	rt.AddHandler(
		handler.NewContextHandler(router.NewRecipientListRouter(container)),
	)
	rt.AddHandler(
		handler.NewContextHandler(handler.NewReplyConsumerHandler()),
	)

	if b.afterInterceptors != nil {
		for _, afterInterceptors := range b.afterInterceptors {
			rt.AddHandler(handler.NewContextHandler(afterInterceptors))
		}
	}

	var messageProcessor message.MessageHandler
	messageProcessor = rt

	if b.deadLetterChannel != "" {
		deadLetterChannel, err := container.Get(b.deadLetterChannel)
		if err != nil {
			panic(fmt.Sprintf("[gateway-builder] [dead-letter] %s", err))
		}
		messageProcessor = handler.NewDeadLetter(
			deadLetterChannel.(message.PublisherChannel),
			messageProcessor,
		)
	}

	return NewGateway(messageProcessor, b.replyChannelName, b.requestChannelName), nil
}

type Gateway struct {
	messageProcessor   message.MessageHandler
	replyChannelName   string
	requestChannelName string
}

func NewGateway(
	messageProcessor message.MessageHandler,
	replyChannelName string,
	requestChannelName string,
) *Gateway {
	return &Gateway{
		messageProcessor:   messageProcessor,
		replyChannelName:   replyChannelName,
		requestChannelName: requestChannelName,
	}
}

func (g *Gateway) Execute(
	parentContext context.Context,
	msg *message.Message,
) (any, error) {
	opCtx, cancel := context.WithCancel(parentContext)
	defer cancel()

	responseChannel := make(chan any)
	go g.executeAsync(opCtx, responseChannel, msg)

	select {
	case result := <-responseChannel:
		switch v := result.(type) {
		case *message.Message:
			return v.GetPayload(), nil
		case error:
			return nil, v
		default:
			return nil, fmt.Errorf("invalid response type")
		}
	case <-opCtx.Done():
		return nil, opCtx.Err()
	}
}

func (g *Gateway) executeAsync(ctx context.Context, responseChannel chan<- any, msg *message.Message) {
	defer close(responseChannel)

	messageToProcess := message.NewMessageBuilderFromMessage(msg)
	messageToProcess.WithChannelName(g.requestChannelName)
	messageToProcess.WithContext(ctx)
	if g.replyChannelName != "" {
		messageToProcess.WithReplyChannelName(g.replyChannelName)
	}

	internalReplyChannel := g.makeInternalChannel()
	messageToProcess.WithReplyChannel(internalReplyChannel)

	resultMessage, err := g.messageProcessor.Handle(ctx, messageToProcess.Build())
	if err != nil {
		internalReplyChannel.Close()
		slog.Error("Failed to process message:",
			"messageId", messageToProcess.Build().GetHeaders().MessageId,
			"reason", err.Error(),
		)
		responseChannel <- err
	}

	select {
	case <-ctx.Done():
		responseChannel <- fmt.Errorf("[gateway]: Context cancelled after processing, before sending result")
		return
	default:
	}

	responseChannel <- resultMessage
}

func (g *Gateway) makeInternalChannel() *channel.PointToPointChannel {
	internalChannelName := uuid.New().String()
	chn := channel.NewPointToPointChannel(internalChannelName)
	return chn
}
