package message

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

const (
	Command MessageType = iota
	Query
	Event
	Document
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
		Receive() (*Message, error)
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
	case Document:
		return "Document"
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
		MessageId     string
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
	messageId := uuid.New().String()
	return &messageHeaders{
		Route:         route,
		MessageType:   messageType,
		Schema:        schema,
		ContentType:   contentType,
		Timestamp:     time.Now(),
		ReplyChannel:  replyChannel,
		CustomHeaders: make(customHeaders),
		CorrelationId: correlationId,
		ChannelName:   channelName,
		MessageId:     messageId,
	}
}

func (m *messageHeaders) SetCustomHeaders(data customHeaders) {
	m.CustomHeaders = data
}

func (m *messageHeaders) MarshalJSON() ([]byte, error) {
	chs, err := json.Marshal(m.CustomHeaders)
	if err != nil {
		panic("[custom-header] cannot marshal.")
	}
	return json.Marshal(struct {
		Route         string    `json:"route"`
		Type          string    `json:"type"`
		Schema        string    `json:"schema"`
		ContentType   string    `json:"contentType"`
		Timestamp     time.Time `json:"timestamp"`
		ReplyChannel  string    `json:"replyChannel"`
		CustomHeaders string    `json:"customHeaders"`
		CorrelationId string    `json:"correlationId"`
		ChannelName   string    `json:"channelName"`
		MessageId     string    `json:"messageId"`
	}{
		m.Route,
		m.MessageType.String(),
		m.Schema,
		m.ContentType,
		m.Timestamp,
		m.ReplyChannel.Name(),
		string(chs),
		m.CorrelationId,
		m.ChannelName,
		m.MessageId,
	})
}

type Message struct {
	payload         any
	headers         *messageHeaders
}

func NewMessage(
	payload any,
	headers *messageHeaders,
) *Message {
	return &Message{
		payload: payload,
		headers: headers,
	}
}

func (m *Message) GetPayload() any {
	return m.payload
}

func (m *Message) GetHeaders() *messageHeaders {
	return m.headers
}

func (m *Message) ReplyRequired() bool {
	return m.headers.MessageType == Command || m.headers.MessageType == Query
}

func (m *Message) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Payload any          `json:"payload"`
		Headers *messageHeaders `json:"headers"`
	}{
		m.payload,
		m.headers,
	})
}
