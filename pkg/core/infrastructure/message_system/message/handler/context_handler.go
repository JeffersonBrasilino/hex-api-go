package handler

import (
	"context"

	"github.com/hex-api-go/pkg/core/infrastructure/message_system/message"
)

type contextHandler struct {
	handler message.MessageHandler
}

func NewContextHandler(handler message.MessageHandler) *contextHandler {
	return &contextHandler{handler: handler}
}

func (h *contextHandler) Handle(ctx context.Context, msg *message.Message) (*message.Message, error) {
	select {
	case <-ctx.Done():
		switch ctx.Err() {
		case context.DeadlineExceeded:
			return nil, context.DeadlineExceeded
		case context.Canceled:
			return nil, context.Canceled
		default:
			return nil, ctx.Err()
		}
	default:
	}
	return h.handler.Handle(ctx, msg)
}
