package bus

import (
	"fmt"
)

type (
	baseHandler interface {
		Trigger() any
	}
	CommandHandler interface {
		baseHandler
		Handle(command any) (any, error)
	}
)

var (
	commandHandlers = map[string]CommandHandler{}
)

func RegisterCommandHandler(action string, executor CommandHandler) {
	if _, exists := commandHandlers[action]; exists {
		panic(fmt.Errorf("handler for action %s already registered", action))
	}
	commandHandlers[action] = executor
}

func GetActionHandler(actionName string) (any, error) {
	handler, exists := commandHandlers[actionName]
	if !exists {
		return nil, fmt.Errorf("handler for action %s has not registered", actionName)
	}
	return handler, nil
}
