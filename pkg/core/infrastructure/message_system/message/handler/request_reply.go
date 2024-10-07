package handler

import "github.com/hex-api-go/pkg/core/infrastructure/message_system/message"

type requestReply struct{}

func NewRequestReply() *requestReply {
	return &requestReply{}
}

func (r *requestReply) Process(msg message.Message) (*message.GenericMessage, error) {

	return nil, nil
}
