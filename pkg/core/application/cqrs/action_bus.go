package cqrs

import (
	"fmt"
	"reflect"
)

var actionsHandlers = map[reflect.Type]any{}

func RegisterActionHandler[TAction any](handler ActionHandler[TAction]) error {
	var action TAction
	actionType := reflect.TypeOf(action)
	_, exists := actionsHandlers[actionType]
	if exists {
		return fmt.Errorf("handler form action %s already registered", actionType)
	}

	actionsHandlers[actionType] = handler
	return nil
}

func Send[TAction any](action TAction) (any, error) {
	handler, exists := actionsHandlers[reflect.TypeOf(action)]
	if !exists {
		return nil, fmt.Errorf("no handler for action %T", action)
	}
	
	handlerAssign, ok := handler.(ActionHandler[TAction])
	
	if !ok {
		return nil,
			fmt.Errorf("handler for action %T not implemented for ActionHandler interface", action)
	}

	return handlerAssign.Handle(action)
}

func Register2(action any, handler any){
	fmt.Println(action, handler)
}