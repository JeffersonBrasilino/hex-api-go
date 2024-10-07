package internalchannel

import (
	"github.com/hex-api-go/pkg/core/infrastructure/message_system/message"
)

type pubSubChannel struct {
	channel chan message.Message
	name    string
}

func NewPubSubChannel(name string) *pubSubChannel {
	return &pubSubChannel{
		name:    name,
		channel: make(chan message.Message),
	}
}

func (p *pubSubChannel) Send(msg message.Message) error {
	p.channel <- msg
	return nil
}

func (p *pubSubChannel) Subscribe(callable func(msg any)) error {
	go func(ch <-chan message.Message) {
		for m := range ch {
			callable(m)
		}
	}(p.channel)
	return nil
}

func (p *pubSubChannel) Receive() (any, error) {
	result := <-p.channel
	return result.(*message.GenericMessage), nil
}

func (p *pubSubChannel) Shutdown() {
	close(p.channel)
}

func (p *pubSubChannel) Name() string {
	return p.name
}

func (p *pubSubChannel) Unsubscribe() error {
	p.Shutdown()
	return nil
}
