package channel

import (
	"fmt"

	"github.com/hex-api-go/pkg/core/infrastructure/message_system/container"
	"github.com/hex-api-go/pkg/core/infrastructure/message_system/message"
)

func DirectChannelReferenceName(name string) string {
	return fmt.Sprintf("direct-channel:%s", name)
}

type DirectChannelBuilder struct {
	channelName string
}

func NewDirectChannelBuilder(channelName string) *DirectChannelBuilder {
	return &DirectChannelBuilder{channelName: channelName}
}

func (c *DirectChannelBuilder) GetChannelName() string {
	return c.channelName
}

func (c *DirectChannelBuilder) Build(
	container container.Container[any, any],
) *DirectChannel {
	if container.Has(DirectChannelReferenceName(c.channelName)) {
		panic(fmt.Sprintf("direct-channel %s already exists", c.channelName))
	}
	return NewDirectChannel(c.channelName)
}

type DirectChannel struct {
	channel chan *message.Message
	name    string
}

func NewDirectChannel(name string) *DirectChannel {
	return &DirectChannel{
		name:    name,
		channel: make(chan *message.Message),
	}
}

func (p *DirectChannel) Send(msg *message.Message) error {
	p.channel <- msg
	return nil
}

func (p *DirectChannel) Subscribe(callable message.MessageHandler) {
	go func(ch <-chan *message.Message) {
		for m := range ch {
			callable.Handle(m)
		}
	}(p.channel)
}

func (p *DirectChannel) Receive() (any, error) {
	result := <-p.channel
	close(p.channel)
	return result, nil
}

func (p *DirectChannel) Unsubscribe() error {
	close(p.channel)
	return nil
}

func (p *DirectChannel) Name() string {
	return p.name
}
