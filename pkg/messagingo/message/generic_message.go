package message

type genericMessage struct {
	Payload []byte          `json:"payload"`
	Headers *messageHeaders `json:"headers"`
}

func CreateGenericMessage(payload []byte, headers *messageHeaders) *genericMessage {
	return &genericMessage{
		Payload: payload,
		Headers: headers,
	}
}

func (m *genericMessage) GetPayload() []byte {
	return m.Payload
}

func (m *genericMessage) GetHeaders() *messageHeaders {
	return m.Headers
}
