package message

import "time"

type customHeaders map[string]string
type messageHeaders struct {
	Route         string        `json:"route"`
	Type          string        `json:"type"` //command|event|message
	Schema        string        `json:"schema"`
	ContentType   string        `json:"contentType"`
	Timestamp     time.Time     `json:"timestamp"`
	ReplyChannel  string        `json:"replyChannel"`
	ErrorChannel  string        `json:"errorChannel"`
	Version       string        `json:"version"`
	CustomHeaders customHeaders `json:"customHeaders"`
	CorrelationId string        `json:"correlationId"`
}

func CreateMessageHeaders() *messageHeaders {
	return &messageHeaders{}
}

func (m *messageHeaders) GetRoute() string {
	return m.Route
}

func (m *messageHeaders) GetSchema() string {
	return m.Schema
}

func (m *messageHeaders) GetContentType() string {
	return m.ContentType
}

func (m *messageHeaders) GetTimestamp() time.Time {
	return m.Timestamp
}

func (m *messageHeaders) GetReplyChannel() string {
	return m.ReplyChannel
}

func (m *messageHeaders) GetErrorChannel() string {
	return m.ErrorChannel
}

func (m *messageHeaders) GetCustomHeaders() customHeaders {
	return m.CustomHeaders
}

func (m *messageHeaders) GetType() string {
	return m.Type
}

func (m *messageHeaders) GetVersion() string {
	return m.Version
}

func (m *messageHeaders) GetCorrelationId() string {
	return m.CorrelationId
}

func (m *messageHeaders) AddCustomHeaders(key, value string) *messageHeaders {
	if m.CustomHeaders == nil {
		m.CustomHeaders = make(map[string]string)
	}
	m.CustomHeaders[key] = value
	return m
}
