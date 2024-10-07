package handler

import (
	"fmt"

	"github.com/hex-api-go/pkg/messagingo/message"
)

type serviceActivator struct {
	outputHandler MessageHandler
}

func NewServiceActivator(outputHandler MessageHandler) *serviceActivator {
	return &serviceActivator{
		outputHandler: outputHandler,
	}
}

func (s *serviceActivator) Handle(message message.Message) (message.Message, error) {
	fmt.Println("service activator handle")
	s.outputHandler.Handle(message)
	return nil, nil
}
