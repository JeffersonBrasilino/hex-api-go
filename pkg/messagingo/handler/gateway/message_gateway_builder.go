package gateway

import (
	"github.com/hex-api-go/pkg/messagingo/container"
	"github.com/hex-api-go/pkg/messagingo/handler"
)

func BuildGateway(
	container container.Container,
	messageHandler handler.MessageHandler,
	afterMiddleware []handler.MessageHandler,
	beforeMiddleware []handler.MessageHandler,
) *defaultMessageGateway {
	messageHandlerPipeline := buildMessageHandlerPipeline(messageHandler, afterMiddleware, beforeMiddleware)
	return CreateDefaultMessageGateway(messageHandlerPipeline)
}

func buildMessageHandlerPipeline(
	messageHandler handler.MessageHandler,
	afterMiddleware []handler.MessageHandler,
	beforeMiddleware []handler.MessageHandler,
) *handler.MessageHandlerPipeline {
	pipeline := handler.NewMessageHandlerPipeline()
	for _, before := range beforeMiddleware {
		pipeline.AddHandler(before)
	}

	pipeline.AddHandler(handler.NewServiceActivator(messageHandler))

	for _, after := range afterMiddleware {
		pipeline.AddHandler(after)
	}

	return pipeline
}
