package bus

import "fmt"

type EventBus interface {
	Publish(event any, route string, aditionalMetadata map[string]string)
}

type eventBus struct{}

func NewEventBus() *eventBus {
	return &eventBus{}
}

func (m *eventBus) Publish(event any, route string, aditionalMetadata map[string]string) {
	fmt.Println("event bus ===> publish")
}
