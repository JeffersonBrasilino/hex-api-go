package channel

import (
	"fmt"
)

type DirectChannelConfiguration struct {
	channelName string
}

func NewDirectChannelConfiguration(channelName string) *DirectChannelConfiguration {
	return &DirectChannelConfiguration{
		channelName: channelName,
	}
}

func (c *DirectChannelConfiguration) GetChannelName() string {
	return c.channelName
}

func (c *DirectChannelConfiguration) ReferenceName() string {
	return fmt.Sprintf("direct-channel: %s", c.channelName)
}