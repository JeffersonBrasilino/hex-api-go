package router

import (
	"context"

	"github.com/hex-api-go/pkg/core/infrastructure/message_system/message"
)

type router struct {
	handlers []message.MessageHandler
}

func NewRouter() *router {
	return &router{
		handlers: []message.MessageHandler{},
	}
}

func (r *router) AddHandler(handler message.MessageHandler) *router {
	r.handlers = append(r.handlers, handler)
	return r
}

func (r *router) Handle(ctx context.Context, msg *message.Message) (*message.Message, error) {
	var resultMessage = msg
	var resulError error
	for _, r := range r.handlers {
		if resultMessage == nil {
			break
		}
		resultMessage, resulError = r.Handle(ctx, resultMessage)
		if resulError != nil {
			break
		}
	}

	return resultMessage, resulError
}
