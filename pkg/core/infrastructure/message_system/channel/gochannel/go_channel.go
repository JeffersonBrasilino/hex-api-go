package gochannel

import (
	"github.com/hex-api-go/pkg/core/infrastructure/message_system/message"
)

type GoChannel struct {
	name    string
	channel chan *message.Message
}

func NewGoChannel(
	name string,
) *GoChannel {
	return &GoChannel{
		name:    name,
		channel: make(chan *message.Message),
	}
}

func (c *GoChannel) Send(msg *message.Message) error {
	c.channel <- msg
	return nil
}

func (c *GoChannel) Subscribe(callable func(msg any)) {
	go func(ch <-chan *message.Message) {
		for m := range ch {
			callable(m)
		}
	}(c.channel)
}

func (c *GoChannel) Receive() *message.Message {
	result := <-c.channel
	return result
}

func (c *GoChannel) Shutdown() {
	close(c.channel)
}

func (c *GoChannel) Name() string {
	return c.name
}
