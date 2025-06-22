package message

import (
	"context"
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
		Handle(ctx context.Context, message *Message) (*Message, error)
	}
	PublisherChannel interface {
		Name() string
		Send(ctx context.Context, message *Message) error
	}
	ConsumerChannel interface {
		Name() string
		Receive() (*Message, error)
		Close() error
	}
	SubscriberChannel interface {
		Name() string
		Subscribe(callable ...func(m *Message))
		Unsubscribe() error
	}
	Gateway interface {
		Execute(parentContext context.Context, msg *Message) (any, error)
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
)

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

func (m *messageHeaders) SetCustomHeaders(data customHeaders) {
	m.CustomHeaders = data
}

func (m *messageHeaders) MarshalJSON() ([]byte, error) {

	chs, err := json.Marshal(m.CustomHeaders)
	if err != nil {
		panic("[custom-header] cannot marshal.")
	}

	var customHeaders string

	if len(chs) > 2 {
		customHeaders = string(chs)
	}

	var replyChannelName = m.ReplyChannelName
	if m.ReplyChannel != nil {
		replyChannelName = m.ReplyChannel.Name()
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
		replyChannelName,
		customHeaders,
		m.CorrelationId,
		m.ChannelName,
		m.MessageId,
	})
}

type Message struct {
	payload any
	headers *messageHeaders
	context context.Context
}

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

func (m *Message) GetPayload() any {
	return m.payload
}

func (m *Message) GetHeaders() *messageHeaders {
	return m.headers
}

func (m *Message) SetContext(ctx context.Context) {
	m.context = ctx
}

func (m *Message) GetContext() context.Context {
	return m.context
}

func (m *Message) ReplyRequired() bool {
	return m.headers.MessageType == Command || m.headers.MessageType == Query
}

func (m *Message) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Payload any             `json:"payload"`
		Headers *messageHeaders `json:"headers"`
	}{
		m.payload,
		m.headers,
	})
}
