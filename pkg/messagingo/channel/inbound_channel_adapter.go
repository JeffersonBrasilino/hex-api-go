package channel

import "fmt"

type inboundChannelAdapter struct {
	messageGateway any
}

func CreateInboundChannelAdapter(
	gateway any,
) *inboundChannelAdapter {
	return &inboundChannelAdapter{
		gateway,
	}
}

func (c *inboundChannelAdapter) ReceiveMessage() ([]byte, error) {
	fmt.Println("DEFAULT CHANNEL >>RECEIVE MESSAGE")
	return []byte("okok"), nil
}
