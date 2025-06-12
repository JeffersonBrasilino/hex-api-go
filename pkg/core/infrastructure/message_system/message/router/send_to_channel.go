package router

import "github.com/hex-api-go/pkg/core/infrastructure/message_system/message"

type sendToChannel struct {
	channel message.PublisherChannel
}

func NewSendToChannel(channel message.PublisherChannel) *sendToChannel {
	return &sendToChannel{
		channel: channel,
	}
}

func (s *sendToChannel) Handle(
	msg *message.Message,
) (*message.Message, error) {
	err := s.channel.Send(msg)
	if err != nil {
		return nil, err
	}
	return msg, nil
}
