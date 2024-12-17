package action

import (
	"fmt"

	"github.com/hex-api-go/pkg/core/infrastructure/message_system/container"
	"github.com/hex-api-go/pkg/core/infrastructure/message_system/message"

	//"github.com/hex-api-go/pkg/core/infrastructure/message_system/message/channel"
	"github.com/hex-api-go/pkg/core/infrastructure/message_system/message/channel"
	"github.com/hex-api-go/pkg/core/infrastructure/message_system/message/endpoint"
)

var (
	actionHandlersToBuild = container.NewGenericContainer[
		string,
		message.MessageHandler,
	]()
)

func AddActionHandler[
	T Action,
	U actionHandlerReturn,
](handler ActionHandler[T, U]) {
	action := *new(T)
	if actionHandlersToBuild.Has(action.Name()) {
		panic(
			fmt.Sprintf(
				"hander for %s already exists",
				action.Name(),
			),
		)
	}

	actionHandlersToBuild.Set(action.Name(), endpoint.NewServiceActivator(
		NewActionActivator(handler),
		"Handle",
	))
}

func Build(container container.Container[any, any]) {
	for action, activator := range actionHandlersToBuild.GetAll() {
		channel := channel.NewPointToPointChannel(action)
		channel.Subscribe(func(msg any) {
			activator.Handle(msg.(*message.Message))
		})
		container.Set(ActionReferenceName(action), channel)
	}
}
