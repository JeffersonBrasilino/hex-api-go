package handler

import (
	"github.com/hex-api-go/pkg/messagingo/message"
)

type MessageHandlerPipeline struct {
	handlers []MessageHandler
}

func NewMessageHandlerPipeline(handlers ...MessageHandler) *MessageHandlerPipeline {
	if handlers == nil {
		handlers = []MessageHandler{}
	}
	return &MessageHandlerPipeline{handlers: handlers}
}

func (h *MessageHandlerPipeline) AddHandler(handler MessageHandler) *MessageHandlerPipeline {
	h.handlers = append(h.handlers, handler)
	return h
}

func (h *MessageHandlerPipeline) Handle(message message.Message) (message.Message, error) {
	messageToProcess := message
	for _, handler := range h.handlers {
		res, err := handler.Handle(messageToProcess)
		if err != nil {
			return nil, err
		}
		messageToProcess = res
	}

	return messageToProcess, nil
}
