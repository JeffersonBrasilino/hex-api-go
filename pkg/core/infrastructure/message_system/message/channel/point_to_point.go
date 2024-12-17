package channel

import (
	"errors"
	"fmt"

	"github.com/hex-api-go/pkg/core/infrastructure/message_system/message"
)

func PointToPointReferenceName(name string) string {
	return fmt.Sprintf("point-to-point-channel:%s", name)
}

type PointToPointChannel struct {
	name    string
	channel chan *message.Message
}

func NewPointToPointChannel(
	name string,
) *PointToPointChannel {
	return &PointToPointChannel{
		name:    name,
		channel: make(chan *message.Message),
	}
}

func (c *PointToPointChannel) Send(msg *message.Message) error {
	go func(ch chan<- *message.Message) {
		c.channel <- msg
	}(c.channel)
	return nil
}

func (c *PointToPointChannel) Subscribe(callable func(m any)) {
	go func(ch <-chan *message.Message) {
		for {
			m, hasOpen := <-ch
			if !hasOpen {
				break
			}
			callable(m)
		}
	}(c.channel)
}

func (c *PointToPointChannel) Receive() (any, error) {
	result, hasOpen := <-c.channel
	if !hasOpen {
		return nil, errors.New("channel has not been closed")
	}
	return result, nil
}

func (c *PointToPointChannel) Close() error {
	close(c.channel)
	return nil
}

func (c *PointToPointChannel) Name() string {
	return fmt.Sprintf("%s", c.name)
}
