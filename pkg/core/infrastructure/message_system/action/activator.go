package action

import (
	"encoding/json"
	"fmt"

	msgHandler "github.com/hex-api-go/pkg/core/infrastructure/message_system/message"
)

type (
	actionHandlerReturn interface {
		any | ~string | ~int | ~bool
	}

	Action interface {
		Type() msgHandler.MessageType
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
	message *msgHandler.Message,
) (*msgHandler.Message, error) {

	var action TInput
	if err := json.Unmarshal(message.GetPayload(), &action); err != nil {
		return nil, err
	}

	var result TOutput
	result, err := c.handler.Handle(action)
	if err != nil {
		return nil, fmt.Errorf("cannot handle action %v: %v", action, err.Error())
	}

	payload, err := json.Marshal(result)
	if err != nil {
		return nil, fmt.Errorf("cannot marshal response: %s", err.Error())
	}

	resultMessage := msgHandler.NewMessageBuilderFromMessage(message).
		WithPayload(payload).
		Build()

	resultMessage.SetInternalPayload(result)

	return resultMessage, nil
}

func ActionReferenceName(name string) string {
	return fmt.Sprintf("action:%s", name)
}
