package router

import (
	"github.com/hex-api-go/pkg/core/infrastructure/message_system/container"
	"github.com/hex-api-go/pkg/core/infrastructure/message_system/message"
)

type messageRouterBuilder struct {
	filters                []message.MessageHandler
	routerComponents       []message.MessageHandler
	messageSystemContainer container.Container[any, any]
}

func NewMessageRouterBuilder(
	container container.Container[any, any],
) *messageRouterBuilder {
	return &messageRouterBuilder{
		filters:                []message.MessageHandler{},
		routerComponents:       []message.MessageHandler{},
		messageSystemContainer: container,
	}
}

func (b *messageRouterBuilder) WithFilter(
	filter FilterFunc,
) *messageRouterBuilder {
	b.filters = append(b.filters, NewMessageFilter(filter))
	return b
}

func (b *messageRouterBuilder) WithRouterComponent(
	component message.MessageHandler,
) *messageRouterBuilder {
	b.routerComponents = append(b.routerComponents, component)
	return b
}

func (b *messageRouterBuilder) WithRecipientListRouter() *messageRouterBuilder {
	recipientListRouter := NewRecipientListRouter(b.messageSystemContainer)
	b.routerComponents = append(b.routerComponents, recipientListRouter)
	return b
}

func (b *messageRouterBuilder) Build() *messageRouter {
	router := NewMessageRouter()
	for _, filter := range b.filters {
		router.Add(filter)
	}

	for _, component := range b.routerComponents {
		router.Add(component)
	}

	return router
}

type messageRouter struct {
	routers []message.MessageHandler
}

func NewMessageRouter() *messageRouter {
	return &messageRouter{
		routers: []message.MessageHandler{},
	}
}

func (r *messageRouter) Add(router message.MessageHandler) *messageRouter {
	r.routers = append(r.routers, router)
	return r
}

func (r *messageRouter) Handle(msg *message.Message) (*message.Message, error) {
	var resultMessage *message.Message
	var resulError error
	for _, r := range r.routers {
		resultMessage, resulError = r.Handle(msg)
		if resulError != nil || resultMessage == nil {
			break
		}
	}

	return resultMessage, resulError
}
