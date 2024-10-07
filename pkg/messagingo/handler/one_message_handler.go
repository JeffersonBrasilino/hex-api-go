package handler

import (
	"fmt"

	"github.com/hex-api-go/pkg/messagingo/message"
)

type OneMessageHandler struct{}

func (handler *OneMessageHandler) Handler(message message.Message) (message.Message, error) {
	fmt.Println("OneMessageHandler called ")
	return message, nil
}

type TwoMessageHandler struct{}

func (handler *TwoMessageHandler) Handler(message message.Message) (message.Message, error) {
	fmt.Println("TwoMessageHandler called ")
	return message, nil
}
