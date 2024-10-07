package messagingo

import "github.com/hex-api-go/pkg/messagingo/message"

type MessageConverter interface {
	ToFormat(message message.Message) any
	ToMessage(message any) *message.Message
}
