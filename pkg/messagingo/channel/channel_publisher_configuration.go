package channel

import (
	"github.com/hex-api-go/pkg/messagingo/client"
)

type ChannelPublisherConfiguration struct {
	queueOrTopicName string
}

func NewChannelPublisherConfiguration() *ChannelPublisherConfiguration {
	return &ChannelPublisherConfiguration{}
}

func (c *ChannelPublisherConfiguration) WithQueueOrTopicName(
	queueOrTopicName string,
) *ChannelPublisherConfiguration {
	c.queueOrTopicName = queueOrTopicName
	return c
}

func (c *ChannelPublisherConfiguration) GetQueueOrTopicName() string {
	return c.queueOrTopicName
}

func (c *ChannelPublisherConfiguration) Build(connection client.ConnectionClient) *outboundChannelAdapter {
	if c.queueOrTopicName == "" {
		panic("the QueueOrTopicName is required")
	}
	return CreateOutboundChannelAdapter(c.queueOrTopicName, connection)
}
