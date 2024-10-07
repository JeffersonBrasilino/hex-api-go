package message

const (
	Command MessageType = iota
	Query
	Event
)

type (
	customHeaders map[string]string
	MessageType   int8
	Message       interface {
		GetHeaders() *messageHeaders
		GetPayload() []byte
	}
	MessageHandler interface {
		Handle(message Message) (any, error)
	}
	MessageProcessor interface {
		Process(message Message) (Message, error)
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
	return "unknown"
}
