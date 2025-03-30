package action

import (
	"fmt"

	"github.com/hex-api-go/pkg/core/infrastructure/message_system/message"
)

type (
	actionHandlerReturn interface {
		any | ~string | ~int | ~bool
	}

	Action interface {
		Type() message.MessageType
		Name() string
	}

	ActionHandler[T Action, U actionHandlerReturn] interface {
		Handle(action T) (U, error)
	}

	ActionActivator[
		THandler ActionHandler[TInput, TOutput],
		TInput Action,
		TOutput actionHandlerReturn,
	] struct {
		handler THandler
	}
)

func NewActionActivator[
	THandler ActionHandler[TInput, TOutput],
	TInput Action,
	TOutput actionHandlerReturn,
](
	handler THandler,
) *ActionActivator[THandler, TInput, TOutput] {
	return &ActionActivator[THandler, TInput, TOutput]{
		handler: handler,
	}
}

func (c *ActionActivator[THandler, TInput, TOutput]) Handle(
	args TInput,
) (TOutput, error) {

	var result TOutput
	result, err := c.handler.Handle(args)
	return result, err
}

func ActionReferenceName(name string) string {
	return fmt.Sprintf("action:%s", name)
}
