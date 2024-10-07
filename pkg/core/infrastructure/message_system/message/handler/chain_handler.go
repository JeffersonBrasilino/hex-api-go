package handler

import "github.com/hex-api-go/pkg/core/infrastructure/message_system/message"

type chainHandler struct {
	handlers []message.MessageProcessor
}

func NewChainHandler(handlers ...message.MessageProcessor) *chainHandler {
	return &chainHandler{
		handlers: handlers,
	}
}

func (h *chainHandler) Process(msg message.Message) (*message.GenericMessage, error) {
	
	return nil, nil
}
