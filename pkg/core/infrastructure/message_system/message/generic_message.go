package message

type GenericMessage struct {
	Payload []byte          `json:"payload"`
	Headers *messageHeaders `json:"headers"`
}

func NewGenericMessage(payload []byte, headers *messageHeaders) *GenericMessage {
	return &GenericMessage{
		Payload: payload,
		Headers: headers,
	}
}

func (m *GenericMessage) GetPayload() []byte {
	return m.Payload
}

func (m *GenericMessage) GetHeaders() *messageHeaders {
	return m.Headers
}