package cqrs

import (
	"fmt"
	"reflect"

	"github.com/hex-api-go/pkg/core/domain"
)

var actionsHandlers = map[reflect.Type]any{}

/*
container que registra as ações(comando e consultas) de forma global.

Tem como objetvo deixar as ações a serem realizadas registradas e independentes
de qual protocolo/abordagem elas vão ser chamadas.

foi usado o mediator pattern para esta abordagem
*/
func RegisterActionHandler[TAction domain.Action[any]](handler ActionHandler[TAction]) error {
	var action TAction
	actionType := reflect.TypeOf(action)
	_, exists := actionsHandlers[actionType]
	if exists {
		return fmt.Errorf("handler form action %s already registered", actionType)
	}

	actionsHandlers[actionType] = handler
	return nil
}

// envia a ação para seu respectivo manipulador.
func Send[TAction domain.Action[any]](action TAction) (any, error) {
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
