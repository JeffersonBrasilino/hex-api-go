package gateway

import (
	"fmt"

	"github.com/hex-api-go/pkg/messagingo/handler"
)

type MessageGateway interface {
	Execute(arguments ...any) any
}

type defaultMessageGateway struct {
	messageHandlerPipeline handler.MessageHandler
}

func CreateDefaultMessageGateway(messageHandlers handler.MessageHandler) *defaultMessageGateway {
	return &defaultMessageGateway{messageHandlers}
}

func (g *defaultMessageGateway) Execute(arguments ...any) any {
	fmt.Println("defaultMessageGateway")
	return nil
}
