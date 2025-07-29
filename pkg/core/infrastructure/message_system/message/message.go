// Package message provides core message types and interfaces for the message system.
//
// This package defines the fundamental message structures, types, and interfaces
// that form the foundation of the messaging system. It implements message types,
// headers, channels, and handlers that enable flexible message processing and
// routing throughout the system.
//
// The Message implementation supports:
// - Multiple message types (Command, Query, Event, Document)
// - Comprehensive header management
// - Context-aware message processing
// - JSON serialization and deserialization
// - Channel interfaces for different messaging patterns
package message

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// MessageType constants define the different types of messages supported by the
// system.
const (
	Command  MessageType = iota // Command messages for actions
	Query                       // Query messages for data retrieval
	Event                       // Event messages for notifications
	Document                    // Document messages for data transfer
)

// MessageType represents the type of a message in the system.
type MessageType int8

// MessageHandler defines the contract for processing messages in the system.
type MessageHandler interface {
	Handle(ctx context.Context, message *Message) (*Message, error)
}

// PublisherChannel defines the contract for channels that can publish messages.
type PublisherChannel interface {
	Name() string
	Send(ctx context.Context, message *Message) error
}

// ConsumerChannel defines the contract for channels that can consume messages.
type ConsumerChannel interface {
	Name() string
	Receive() (*Message, error)
	Close() error
}

// SubscriberChannel defines the contract for channels that support subscription
// patterns.
type SubscriberChannel interface {
	Name() string
	Subscribe(callable ...func(m *Message))
	Unsubscribe() error
}

// InboundChannelAdapter defines the contract for inbound channel adapters that
// receive messages from external sources.
type InboundChannelAdapter interface {
	ReferenceName() string
	DeadLetterChannelName() string
	AfterProcessors() []MessageHandler
	BeforeProcessors() []MessageHandler
	ReceiveMessage(ctx context.Context) (*Message, error)
	Close() error
}

// String returns the string representation of a MessageType.
//
// Returns:
//   - string: the string representation of the message type
func (m MessageType) String() string {
	switch m {
	case Command:
		return "Command"
	case Query:
		return "Query"
	case Event:
		return "Event"
	case Document:
		return "Document"
	}
	return "Message"
}

// customHeaders represents a map of custom header key-value pairs.
type customHeaders map[string]string

// messageHeaders contains all the metadata associated with a message including
// routing information, timestamps, and custom headers.
type messageHeaders struct {
	Route            string
	MessageType      MessageType
	Timestamp        time.Time
	ReplyChannel     PublisherChannel
	CustomHeaders    customHeaders
	CorrelationId    string
	ChannelName      string
	MessageId        string
	ReplyChannelName string
}

// NewMessageHeaders creates a new message headers instance with the specified
// parameters and automatically generated message ID and timestamp.
//
// Parameters:
//   - route: the routing information for the message
//   - messageType: the type of the message
//   - replyChannel: the channel for reply messages
//   - correlationId: the correlation identifier for message tracking
//   - channelName: the name of the target channel
//   - replyChannelName: the name of the reply channel
//
// Returns:
//   - *messageHeaders: new message headers instance
func NewMessageHeaders(
	route string,
	messageType MessageType,
	replyChannel PublisherChannel,
	correlationId string,
	channelName string,
	replyChannelName string,
) *messageHeaders {
	messageId := uuid.New().String()
	return &messageHeaders{
		Route:            route,
		MessageType:      messageType,
		Timestamp:        time.Now(),
		ReplyChannel:     replyChannel,
		CustomHeaders:    make(customHeaders),
		CorrelationId:    correlationId,
		ChannelName:      channelName,
		MessageId:        messageId,
		ReplyChannelName: replyChannelName,
	}
}

// SetCustomHeaders sets the custom headers for the message.
//
// Parameters:
//   - data: the custom headers to be set
func (m *messageHeaders) SetCustomHeaders(data customHeaders) {
	m.CustomHeaders = data
}

// MarshalJSON implements the json.Marshaler interface for message headers.
//
// Returns:
//   - []byte: the JSON representation of the headers
//   - error: error if marshaling fails
func (m *messageHeaders) MarshalJSON() ([]byte, error) {

	chs, err := json.Marshal(m.CustomHeaders)
	if err != nil {
		panic("[custom-header] cannot marshal.")
	}

	var customHeaders string

	if len(chs) > 2 {
		customHeaders = string(chs)
	}

	return json.Marshal(struct {
		Route         string    `json:"route"`
		Type          string    `json:"type"`
		Timestamp     time.Time `json:"timestamp"`
		ReplyChannel  string    `json:"replyChannel"`
		CustomHeaders string    `json:"customHeaders"`
		CorrelationId string    `json:"correlationId"`
		ChannelName   string    `json:"channelName"`
		MessageId     string    `json:"messageId"`
	}{
		m.Route,
		m.MessageType.String(),
		m.Timestamp,
		m.ReplyChannelName,
		customHeaders,
		m.CorrelationId,
		m.ChannelName,
		m.MessageId,
	})
}

// Message represents a message in the system with payload, headers, and context.
type Message struct {
	payload any
	headers *messageHeaders
	context context.Context
}

// NewMessage creates a new message instance with the specified payload, headers,
// and context.
//
// Parameters:
//   - payload: the data carried by the message
//   - headers: the message headers containing metadata
//   - context: the context for the message
//
// Returns:
//   - *Message: new message instance
func NewMessage(
	payload any,
	headers *messageHeaders,
	context context.Context,
) *Message {
	return &Message{
		payload: payload,
		headers: headers,
		context: context,
	}
}

// GetPayload returns the payload of the message.
//
// Returns:
//   - any: the message payload
func (m *Message) GetPayload() any {
	return m.payload
}

// GetHeaders returns the headers of the message.
//
// Returns:
//   - *messageHeaders: the message headers
func (m *Message) GetHeaders() *messageHeaders {
	return m.headers
}

// SetContext sets the context for the message.
//
// Parameters:
//   - ctx: the context to be set
func (m *Message) SetContext(ctx context.Context) {
	m.context = ctx
}

// GetContext returns the context of the message.
//
// Returns:
//   - context.Context: the message context
func (m *Message) GetContext() context.Context {
	return m.context
}

// ReplyRequired determines if the message requires a reply based on its type.
// Commands and Queries typically require replies, while Events and Documents
// do not.
//
// Returns:
//   - bool: true if the message requires a reply, false otherwise
func (m *Message) ReplyRequired() bool {
	return m.headers.MessageType == Command || m.headers.MessageType == Query
}

// MarshalJSON implements the json.Marshaler interface for messages.
//
// Returns:
//   - []byte: the JSON representation of the message
//   - error: error if marshaling fails
func (m *Message) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Payload any             `json:"payload"`
		Headers *messageHeaders `json:"headers"`
	}{
		m.payload,
		m.headers,
	})
}
