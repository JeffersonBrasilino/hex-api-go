package handler

import (
	"encoding/json"
	"fmt"

	"github.com/hex-api-go/pkg/core/infrastructure/message_system/bus"
	msgHandler "github.com/hex-api-go/pkg/core/infrastructure/message_system/message"
)

type ServiceActivator struct {
	service any
}

func NewServiceActivator(service any) *ServiceActivator {
	return &ServiceActivator{
		service: service,
	}
}

func (c *ServiceActivator) Handle(message msgHandler.Message) (any, error) {
	switch c.service.(type) {
	case bus.CommandHandler:
		return c.executeCommand(message.GetPayload())
	case msgHandler.MessageHandler:
		return c.executeRawHandler(message)
	default:
		return nil, fmt.Errorf(fmt.Sprintf("has no handler for %s", message.GetHeaders().GetRoute()))
	}
}

func (c *ServiceActivator) executeCommand(data []byte) (any, error) {
	commandHandler := c.service.(bus.CommandHandler)
	command := commandHandler.Trigger()
	if err := json.Unmarshal(data, &command); err != nil {
		return nil, err
	}
	commandHandler.Handle(command)
	return nil, nil
}

func (c *ServiceActivator) executeRawHandler(data msgHandler.Message) (any, error) {
	rawHandler := c.service.(msgHandler.MessageHandler)
	return rawHandler.Handle(data)
}
