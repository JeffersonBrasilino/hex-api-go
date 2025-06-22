package handler

import (
	"context"
	"log/slog"

	"github.com/hex-api-go/pkg/core/infrastructure/message_system/message"
)

type deadLetter struct {
	channel message.PublisherChannel
	handler message.MessageHandler
}

func NewDeadLetter(
	channel message.PublisherChannel,
	handler message.MessageHandler,
) *deadLetter {
	return &deadLetter{
		channel: channel,
		handler: handler,
	}
}

func (s *deadLetter) Handle(
	ctx context.Context,
	msg *message.Message,
) (*message.Message, error) {

	_, err := s.handler.Handle(ctx, msg)
	if err != nil {
		slog.Error("[dead-letter] Sending message to dead letter",
			"messageId", msg.GetHeaders().MessageId,
			"reason", err.Error(),
		)

		s.channel.Send(ctx, msg)
	}

	return msg, err
}
