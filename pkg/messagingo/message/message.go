package message

type Message interface {
	GetHeaders() *messageHeaders
	GetPayload() []byte
}
