// Package handler provides message handling components for the message system.
//
// This package implements various message handlers that process and route messages
// through the system. It provides specialized handlers for different message
// processing scenarios including dead letter handling, context management, and
// error handling patterns.
//
// The DeadLetter implementation supports:
// - Failed message handling and routing
// - Dead letter channel integration
// - Error logging and monitoring
// - Graceful error recovery patterns
package handler

import (
	"context"
	"log/slog"

	"github.com/hex-api-go/pkg/core/infrastructure/messagesystem/message"
)

// deadLetter implements the Dead Letter Channel pattern, routing failed messages
// to a designated dead letter channel for further processing or analysis.
type deadLetter struct {
	channel message.PublisherChannel
	handler message.MessageHandler
}

// NewDeadLetter creates a new dead letter handler instance that routes failed
// messages to the specified dead letter channel.
//
// Parameters:
//   - channel: the publisher channel for sending failed messages
//   - handler: the message handler to attempt processing with
//
// Returns:
//   - *deadLetter: configured dead letter handler
func NewDeadLetter(
	channel message.PublisherChannel,
	handler message.MessageHandler,
) *deadLetter {
	return &deadLetter{
		channel: channel,
		handler: handler,
	}
}

// Handle processes a message by attempting to process it with the wrapped handler.
// If processing fails, the message is sent to the dead letter channel for further
// analysis or processing.
//
// Parameters:
//   - ctx: context for timeout/cancellation control
//   - msg: the message to be processed
//
// Returns:
//   - *message.Message: the original message (regardless of processing success)
//   - error: error if processing fails (message is sent to dead letter channel)
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
