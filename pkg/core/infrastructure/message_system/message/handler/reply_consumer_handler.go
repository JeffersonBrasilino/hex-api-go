package handler

import (
	"context"
	"fmt"

	"github.com/hex-api-go/pkg/core/infrastructure/message_system/message"
)

type replyConsumerHandler struct{}

func NewReplyConsumerHandler() *replyConsumerHandler {
	return &replyConsumerHandler{}
}

func (s *replyConsumerHandler) Handle(
	ctx context.Context,
	msg *message.Message,
) (*message.Message, error) {
	replyChannel, ok := msg.GetHeaders().ReplyChannel.(message.ConsumerChannel)

	if !ok {
		return nil, fmt.Errorf("reply channel is not a consumer channel")
	}

	if replyChannel == nil {
		return nil, fmt.Errorf("reply channel not found")
	}

	replyMessage, err := replyChannel.Receive()

	if err != nil {
		return nil, err
	}

	if errorMessage, ok := replyMessage.GetPayload().(error); ok {
		return nil, errorMessage
	}

	return replyMessage, nil
}
