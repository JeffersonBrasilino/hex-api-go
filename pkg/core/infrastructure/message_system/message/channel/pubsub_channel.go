package channel

import (
	"github.com/hex-api-go/pkg/core/infrastructure/message_system/message"
)

type PubSubChannel struct {
	channel chan *message.Message
	name    string
}

func NewPubSubChannel(name string) *PubSubChannel {
	return &PubSubChannel{
		name:    name,
		channel: make(chan *message.Message),
	}
}

func (p *PubSubChannel) Send(msg *message.Message) error {
	p.channel <- msg
	return nil
}

func (p *PubSubChannel) Subscribe(callable ...func(m *message.Message)) {
	go func(ch <-chan *message.Message) {
		for {
			m, hasOpen := <-ch
			if !hasOpen {
				break
			}

			for _, call := range callable {
				call(m)
			}
		}
	}(p.channel)
}

func (p *PubSubChannel) Unsubscribe() error {
	close(p.channel)
	return nil
}
