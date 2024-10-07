package channel

import (
	"fmt"
	"github.com/hex-api-go/pkg/messagingo/client"
	"github.com/hex-api-go/pkg/messagingo/message"
)

type outboundChannelAdapter struct {
	name   string
	client client.ConnectionClient
}

func CreateOutboundChannelAdapter(
	name string,
	client client.ConnectionClient,
) *outboundChannelAdapter {
	return &outboundChannelAdapter{
		name,
		client,
	}
}

func (c *outboundChannelAdapter) Handle(message message.Message) (message.Message, error) {
	fmt.Println("DEFAULT CHANNEL >>RECEIVE MESSAGE")
	return nil, nil
}
