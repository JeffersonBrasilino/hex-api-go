package message

import (
	"encoding/json"
	"time"
)

type messageHeaders struct {
	route                string
	messageType          string
	schema               string
	contentType          string
	timestamp            time.Time
	replyChannel         string
	errorChannel         string
	version              string
	customHeaders        customHeaders
	correlationId        string
	channelName          string
}

func NewMessageHeaders(
	Route,
	Type,
	Schema,
	ContentType,
	ReplyChannel,
	ErrorChannel,
	Version,
	CorrelationId,
	ChannelName string,
) *messageHeaders {
	return &messageHeaders{
		Route,
		Type,
		Schema,
		ContentType,
		time.Now(),
		ReplyChannel,
		ErrorChannel,
		Version,
		make(customHeaders),
		CorrelationId,
		ChannelName,
	}
}

func (m *messageHeaders) SetRoute(value string) {
	m.route = value
}

func (m *messageHeaders) SetType(value MessageType) {
	m.messageType = value.String()
}

func (m *messageHeaders) SetSchema(value string) {
	m.schema = value
}

func (m *messageHeaders) SetContentType(value string) {
	m.contentType = value
}

func (m *messageHeaders) SetReplyChannel(value string) {
	m.replyChannel = value
}

func (m *messageHeaders) SetErrorChannel(value string) {
	m.errorChannel = value
}

func (m *messageHeaders) SetVersion(value string) {
	m.version = value
}

func (m *messageHeaders) SetCorrelationId(value string) {
	m.correlationId = value
}

func (m *messageHeaders) SetChannelName(value string) {
	m.channelName = value
}

func (m *messageHeaders) AddCustomHeaders(key, value string) *messageHeaders {
	m.customHeaders[key] = value
	return m
}

func (m *messageHeaders) SetCustomHeaders(value customHeaders) *messageHeaders {
	m.customHeaders = value
	return m
}

func (m *messageHeaders) GetRoute() string {
	return m.route
}

func (m *messageHeaders) GetReplyChannel() string {
	return m.replyChannel
}

func (m *messageHeaders) GetChannelName() string {
	return m.channelName
}

func (m *messageHeaders) MarshallJSON() ([]byte, error) {
	marshalled, err := json.Marshal(struct {
		Route         string        `json:"route"`
		Type          string        `json:"type"`
		Schema        string        `json:"schema"`
		ContentType   string        `json:"contentType"`
		Timestamp     time.Time     `json:"timestamp"`
		ReplyChannel  string        `json:"replyChannel"`
		ErrorChannel  string        `json:"errorChannel"`
		Version       string        `json:"version"`
		CustomHeaders customHeaders `json:"customHeaders"`
		CorrelationId string        `json:"correlationId"`
		ChannelName   string        `json:"channelName"`
	}{
		m.route,
		m.messageType,
		m.schema,
		m.contentType,
		m.timestamp,
		m.replyChannel,
		m.errorChannel,
		m.version,
		m.customHeaders,
		m.correlationId,
		m.channelName,
	})
	if err != nil {
		return nil, err
	}

	return marshalled, nil
}
