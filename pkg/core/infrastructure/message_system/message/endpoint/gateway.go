// Package endpoint provides gateway functionality for message processing and routing.
//
// This package implements the Gateway pattern from Enterprise Integration Patterns,
// serving as an entry point for message processing with support for interceptors,
// dead letter channels, and reply channels. It provides a centralized message
// processing pipeline with configurable routing and error handling.
//
// The Gateway implementation supports:
// - Message processing with before/after interceptors
// - Dead letter channel integration for failed messages
// - Reply channel support for request-response patterns
// - Asynchronous message processing with context support
// - Configurable routing through recipient list routers
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

// GatewayReferenceName generates a standardized reference name for gateways.
//
// Parameters:
//   - referenceName: the base name for the gateway
//
// Returns:
//   - string: the formatted reference name with prefix
func GatewayReferenceName(referenceName string) string {
	return fmt.Sprintf("gateway:%s", referenceName)
}

// gatewayBuilder provides a fluent interface for configuring gateway instances
// with various options like interceptors, dead letter channels, and reply channels.
type gatewayBuilder struct {
	referenceName      string
	requestChannelName string
	beforeInterceptors []message.MessageHandler
	afterInterceptors  []message.MessageHandler
	deadLetterChannel  string
	replyChannelName   string
}

// Gateway represents a message processing gateway that handles message routing,
// interceptors, and error handling through a configurable processing pipeline.
type Gateway struct {
	messageProcessor   message.MessageHandler
	replyChannelName   string
	requestChannelName string
}

// NewGatewayBuilder creates a new gateway builder instance.
//
// Parameters:
//   - referenceName: unique identifier for the gateway
//   - requestChannelName: name of the channel to process messages from
//
// Returns:
//   - *gatewayBuilder: configured builder instance
func NewGatewayBuilder(
	referenceName string,
	requestChannelName string,
) *gatewayBuilder {
	return &gatewayBuilder{
		referenceName:      referenceName,
		requestChannelName: requestChannelName,
	}
}

// NewGateway creates a new gateway instance.
//
// Parameters:
//   - messageProcessor: the message handler for processing messages
//   - replyChannelName: name of the reply channel
//   - requestChannelName: name of the request channel
//
// Returns:
//   - *Gateway: configured gateway instance
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

// ReferenceName returns the standardized reference name for the gateway.
//
// Returns:
//   - string: the formatted reference name
func (b *gatewayBuilder) ReferenceName() string {
	return GatewayReferenceName(b.referenceName)
}

// WithBeforeInterceptors adds interceptors to be executed before message processing.
//
// Parameters:
//   - interceptors: variable number of message handlers to execute before processing
//
// Returns:
//   - *gatewayBuilder: builder instance for method chaining
func (b *gatewayBuilder) WithBeforeInterceptors(
	interceptors ...message.MessageHandler,
) *gatewayBuilder {
	b.beforeInterceptors = append(b.beforeInterceptors, interceptors...)
	return b
}

// WithAfterInterceptors adds interceptors to be executed after message processing.
//
// Parameters:
//   - interceptors: variable number of message handlers to execute after processing
//
// Returns:
//   - *gatewayBuilder: builder instance for method chaining
func (b *gatewayBuilder) WithAfterInterceptors(
	interceptors ...message.MessageHandler,
) *gatewayBuilder {
	b.afterInterceptors = append(b.afterInterceptors, interceptors...)
	return b
}

// WithDeadLetterChannel sets the dead letter channel for failed messages.
//
// Parameters:
//   - channelName: name of the dead letter channel
//
// Returns:
//   - *gatewayBuilder: builder instance for method chaining
func (b *gatewayBuilder) WithDeadLetterChannel(channelName string) *gatewayBuilder {
	b.deadLetterChannel = channelName
	return b
}

// WithReplyChannel sets the reply channel for request-response patterns.
//
// Parameters:
//   - channelName: name of the reply channel
//
// Returns:
//   - *gatewayBuilder: builder instance for method chaining
func (b *gatewayBuilder) WithReplyChannel(channelName string) *gatewayBuilder {
	b.replyChannelName = channelName
	return b
}

// Build constructs a Gateway from the dependency container with configured
// interceptors, dead letter channel, and reply channel.
//
// Parameters:
//   - container: dependency container containing required components
//
// Returns:
//   - *Gateway: configured gateway instance
//   - error: error if construction fails
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

// Execute processes a message through the gateway's processing pipeline with
// context support and response handling.
//
// Parameters:
//   - parentContext: parent context for timeout/cancellation control
//   - msg: the message to be processed
//
// Returns:
//   - any: the processing result
//   - error: error if processing fails or context is cancelled
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

// executeAsync processes a message asynchronously and sends the result to the
// response channel.
//
// Parameters:
//   - ctx: context for timeout/cancellation control
//   - responseChannel: channel to send processing results
//   - msg: the message to be processed
func (g *Gateway) executeAsync(
	ctx context.Context,
	responseChannel chan<- any,
	msg *message.Message,
) {
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
		responseChannel <- fmt.Errorf(
			"[gateway]: Context cancelled after processing, before sending result",
		)
		return
	default:
	}

	responseChannel <- resultMessage
}

// makeInternalChannel creates an internal point-to-point channel for handling
// reply messages during processing.
//
// Returns:
//   - *channel.PointToPointChannel: internal channel for reply handling
func (g *Gateway) makeInternalChannel() *channel.PointToPointChannel {
	internalChannelName := uuid.New().String()
	chn := channel.NewPointToPointChannel(internalChannelName)
	return chn
}
