package message

import (
	"encoding/json"
	"time"
)

const (
	Command MessageType = iota
	Query
	Event
)

type (
	MessageType    int8
	MessageHandler interface {
		Handle(message *Message) (*Message, error)
	}
	Channel interface {
		Name() string
	}
	PublisherChannel interface {
		Channel
		Send(message *Message) error
	}
	ConsumerChannel interface {
		Channel
		Receive() (any, error)
		Close() error
	}
	SubscriberChannel interface {
		Channel
		Subscribe(callable ...func(m *Message))
		Unsubscribe() error
	}
)

func (m MessageType) String() string {
	switch m {
	case Command:
		return "Command"
	case Query:
		return "Query"
	case Event:
		return "Event"
	}
	return "Message"
}

type (
	customHeaders  map[string]string
	messageHeaders struct {
		Route         string
		MessageType   MessageType
		Schema        string
		ContentType   string
		Timestamp     time.Time
		ReplyChannel  PublisherChannel
		CustomHeaders customHeaders
		CorrelationId string
		ChannelName   string
	}
)

func NewMessageHeaders(
	route string,
	messageType MessageType,
	schema string,
	contentType string,
	replyChannel PublisherChannel,
	correlationId string,
	channelName string,
) *messageHeaders {
	return &messageHeaders{
		route,
		messageType,
		schema,
		contentType,
		time.Now(),
		replyChannel,
		make(customHeaders),
		correlationId,
		channelName,
	}
}

func (m *messageHeaders) SetCustomHeaders(data customHeaders) {
	m.CustomHeaders = data
}

func (m *messageHeaders) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Route         string        `json:"route"`
		Type          string        `json:"type"`
		Schema        string        `json:"schema"`
		ContentType   string        `json:"contentType"`
		Timestamp     time.Time     `json:"timestamp"`
		ReplyChannel  string        `json:"replyChannel"`
		CustomHeaders customHeaders `json:"customHeaders"`
		CorrelationId string        `json:"correlationId"`
		ChannelName   string        `json:"channelName"`
	}{
		m.Route,
		m.MessageType.String(),
		m.Schema,
		m.ContentType,
		m.Timestamp,
		m.ReplyChannel.Name(),
		m.CustomHeaders,
		m.CorrelationId,
		m.ChannelName,
	})
}

type Message struct {
	payload         []byte
	headers         *messageHeaders
	internalPayload any
}

func NewMessage(
	payload []byte,
	headers *messageHeaders,
) *Message {
	return &Message{
		payload: payload,
		headers: headers,
	}
}

func (m *Message) GetPayload() []byte {
	return m.payload
}

func (m *Message) GetHeaders() *messageHeaders {
	return m.headers
}

func (m *Message) SetInternalPayload(instance any) {
	m.internalPayload = instance
}

func (m *Message) GetInternalPayload() any {
	return m.internalPayload
}

func (m *Message) ReplyRequired() bool {
	return m.headers.MessageType == Command || m.headers.MessageType == Query
}

func (m *Message) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Payload []byte          `json:"payload"`
		Headers *messageHeaders `json:"headers"`
	}{
		m.payload,
		m.headers,
	})
}
