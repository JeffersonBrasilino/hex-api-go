package getuser

import "github.com/hex-api-go/pkg/core/infrastructure/message_system/message"

type Query struct {
	DataSource    string
	mapMockedData []float64
}

func NewQuery() *Query {
	return &Query{}
}

func (c *Query) Type() message.MessageType {
	return message.Query
}

func (c *Query) Name() string {
	return "getUser"
}
