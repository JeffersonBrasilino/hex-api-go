package router

import (
	"github.com/hex-api-go/pkg/core/infrastructure/message_system/message"
)

type FilterFunc func(message.Message) bool
type messageFilter struct {
	filterFunc FilterFunc
}

func NewMessageFilter(filterFunc FilterFunc) *messageFilter {
	return &messageFilter{filterFunc: filterFunc}
}

func (f *messageFilter) Handle(msg *message.Message) (*message.Message, error) {
	filterResult := f.filterFunc(*msg)
	if filterResult {
		return msg, nil
	}
	
	return nil, nil
}
