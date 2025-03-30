package router

import (
	"github.com/hex-api-go/pkg/core/infrastructure/message_system/container"
	"github.com/hex-api-go/pkg/core/infrastructure/message_system/message"
)

type MessageRouterBuilder struct {
	filters          []message.MessageHandler
	routerComponents []message.MessageHandler
	hasRecipientList bool
}

func NewMessageRouterBuilder() *MessageRouterBuilder {
	return &MessageRouterBuilder{
		filters:          []message.MessageHandler{},
		routerComponents: []message.MessageHandler{},
	}
}

func (b *MessageRouterBuilder) WithFilter(
	filter FilterFunc,
) *MessageRouterBuilder {
	b.filters = append(b.filters, NewMessageFilter(filter))
	return b
}

func (b *MessageRouterBuilder) WithRouterComponent(
	component message.MessageHandler,
) *MessageRouterBuilder {
	b.routerComponents = append(b.routerComponents, component)
	return b
}

func (b *MessageRouterBuilder) WithRecipientListRouter() *MessageRouterBuilder {
	b.hasRecipientList = true
	return b
}

func (b *MessageRouterBuilder) Build(container container.Container[any, any]) *messageRouter {
	if b.hasRecipientList {
		recipientListRouter := NewRecipientListRouter(container)
		b.routerComponents = append(b.routerComponents, recipientListRouter)
	}

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
