package handler

import "github.com/hex-api-go/pkg/messagingo/message"

type MessageHandler interface {
	Handle(message message.Message) (message.Message, error)
}
