package channel

import (
	"context"
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
	hasOpen bool
}

func NewPointToPointChannel(
	name string,
) *PointToPointChannel {
	return &PointToPointChannel{
		name:    name,
		channel: make(chan *message.Message),
		hasOpen: true,
	}
}

func (c *PointToPointChannel) Send(ctx context.Context, msg *message.Message) error {
	if !c.hasOpen {
		return errors.New("channel has not been opened")
	}

	c.channel <- msg
	return nil
}

func (c *PointToPointChannel) Subscribe(callable func(m *message.Message)) {
	go func(ch <-chan *message.Message) {
		for {
			m, hasOpen := <-ch
			if !hasOpen {
				c.hasOpen = false
				break
			}
			go callable(m)
		}
	}(c.channel)
}

func (c *PointToPointChannel) Receive() (*message.Message, error) {
	result, hasOpen := <-c.channel
	if !hasOpen {
		c.hasOpen = false
		return nil, errors.New("channel has not been opened")
	}
	c.Close()
	return result, nil
}

func (c *PointToPointChannel) Close() error {
	if !c.hasOpen {
		return nil
	}
	c.hasOpen = false
	close(c.channel)
	return nil
}

func (c *PointToPointChannel) Name() string {
	return c.name
}
