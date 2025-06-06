package action

import (
	"fmt"

	"github.com/hex-api-go/pkg/core/infrastructure/message_system/container"
	"github.com/hex-api-go/pkg/core/infrastructure/message_system/message"
	"github.com/hex-api-go/pkg/core/infrastructure/message_system/message/channel"
	"github.com/hex-api-go/pkg/core/infrastructure/message_system/message/endpoint"
)

var (
	actionHandlersToBuild = container.NewGenericContainer[
		string,
		message.MessageHandler,
	]()
)

func ActionReferenceName(name string) string {
	return fmt.Sprintf("action:%s", name)
}

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

// TODO: refatorar e adicionar a função de build no action activator.
func Build(container container.Container[any, any]) {
	fmt.Println("build actions...")
	for action, activator := range actionHandlersToBuild.GetAll() {
		channel := channel.NewPointToPointChannel(action)
		channel.Subscribe(func(msg *message.Message) {
			activator.Handle(msg)
		})
		container.Set(ActionReferenceName(action), channel)
	}
}
