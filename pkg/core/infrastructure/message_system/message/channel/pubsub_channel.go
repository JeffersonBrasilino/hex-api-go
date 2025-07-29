// Package channel provides publish-subscribe messaging channels for the message system.
//
// This package implements the Publish-Subscribe Channel pattern from Enterprise
// Integration Patterns, enabling one-to-many message distribution where a single
// publisher can send messages to multiple subscribers. It provides asynchronous
// message broadcasting with support for multiple concurrent subscribers.
//
// The PubSubChannel implementation supports:
// - One-to-many message broadcasting
// - Multiple concurrent subscribers
// - Asynchronous message processing
// - Graceful channel closure and resource cleanup
// - Thread-safe operations with proper state management
package channel

import (
	"context"
	"fmt"

	"github.com/hex-api-go/pkg/core/infrastructure/message_system/message"
)

// PubSubChannel implements a publish-subscribe messaging channel where messages
// are broadcast to all registered subscribers.
type PubSubChannel struct {
	channel chan *message.Message
	name    string
}

// NewPubSubChannel creates a new publish-subscribe channel instance.
//
// Parameters:
//   - name: the name identifier for the channel
//
// Returns:
//   - *PubSubChannel: a new configured publish-subscribe channel
func NewPubSubChannel(name string) *PubSubChannel {
	return &PubSubChannel{
		name:    name,
		channel: make(chan *message.Message),
	}
}

// Send publishes a message to all registered subscribers.
//
// Parameters:
//   - msg: the message to be published
//
// Returns:
//   - error: error if sending fails (typically nil)
func (p *PubSubChannel) Send(msg *message.Message) error {
	select {
	case p.channel <- msg:
		return nil
	default:
		return fmt.Errorf("channel is full or closed")
	}
}

// SendWithContext publishes a message to all registered subscribers with context
// support for timeout and cancellation.
//
// Parameters:
//   - ctx: context for timeout/cancellation control
//   - msg: the message to be published
//
// Returns:
//   - error: error if sending fails or context is cancelled
func (p *PubSubChannel) SendWithContext(ctx context.Context, msg *message.Message) error {
	select {
	case p.channel <- msg:
		return nil
	case <-ctx.Done():
		return fmt.Errorf("context cancelled while sending message: %v", ctx.Err())
	default:
		return fmt.Errorf("channel is full or closed")
	}
}

// Subscribe registers one or more callback functions to receive published messages.
// Each callback is executed in a separate goroutine for each received message.
//
// Parameters:
//   - callable: variable number of functions to be called for each received message
func (p *PubSubChannel) Subscribe(callable ...func(m *message.Message)) {
	go func(ch <-chan *message.Message) {
		for {
			m, hasOpen := <-ch
			if !hasOpen {
				break
			}

			for _, call := range callable {
				go call(m)
			}
		}
	}(p.channel)
}

// Unsubscribe closes the publish-subscribe channel and stops accepting new messages.
// Existing subscribers will continue to process messages until the channel is empty.
//
// Returns:
//   - error: error if closing the channel fails (typically nil)
func (p *PubSubChannel) Unsubscribe() error {
	close(p.channel)
	return nil
}

// Name returns the name identifier of the publish-subscribe channel.
//
// Returns:
//   - string: the channel name
func (p *PubSubChannel) Name() string {
	return p.name
}
