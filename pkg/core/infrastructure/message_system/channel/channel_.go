package channel

import (
	"fmt"

	"github.com/hex-api-go/pkg/core/infrastructure/message_system/container"
	"github.com/hex-api-go/pkg/core/infrastructure/message_system/message"
)

type (
	Channel interface {
		Name() string
	}
	MessageOutboundChannel interface {
		Channel
		Send(payload message.Message) error
	}
	MessageInboundChannel interface {
		Channel
		Subscribe(callable func(msg any)) error
		Receive() (any, error)
		Unsubscribe() error
	}
	ChannelAdapterBuilder interface {
		Build() any
	}
)

var (
	channelBuilders = []ChannelAdapterBuilder{}
	outboundChannel = container.NewGenericContainer[string, MessageOutboundChannel]()
	inboundChannel  = container.NewGenericContainer[string, MessageInboundChannel]()
)

func AddChannelBuilder(outboundChannel ChannelAdapterBuilder) {
	channelBuilders = append(channelBuilders, outboundChannel)
}

func Build() {
	for _, channel := range channelBuilders {
		buildedChannel := channel.Build()
		if channelConverted, ok := buildedChannel.(MessageOutboundChannel); ok {
			AddOutboundChannel(channelConverted)
		}
		if channelConverted, ok := buildedChannel.(MessageInboundChannel); ok {
			AddInboundChannel(channelConverted)
		}
	}
}

func GetOutboundChannel(channelName string) (MessageOutboundChannel, error) { //TODO: ajustar retorno de interface
	channel, err := outboundChannel.Get(channelName)
	if err != nil {
		return nil, fmt.Errorf("outbound channel %s not found", channelName)
	}
	return channel, nil
}

func AddOutboundChannel(channel MessageOutboundChannel) {
	if outboundChannel.Has(channel.Name()) {
		panic(
			fmt.Sprintf(
				"outbound channel %s already exists",
				channel.Name(),
			),
		)
	}
	outboundChannel.Set(channel.Name(), channel)
}

func AddInboundChannel(channel MessageInboundChannel) {
	if inboundChannel.Has(channel.Name()) {
		panic(
			fmt.Sprintf(
				"inbound channel %s already exists",
				channel.Name(),
			),
		)
	}
	inboundChannel.Set(channel.Name(), channel)
}

func GetInboundChannel(channelName string) (MessageInboundChannel, error) { //TODO: ajustar retorno de interface
	channel, err := inboundChannel.Get(channelName)
	if err != nil {
		return nil, fmt.Errorf("outbound channel %s not found", channelName)
	}
	return channel, nil
}

func RemoveOutboundChannel(channelName string) error {
	return outboundChannel.Remove(channelName)
}

func RemoveInboundChannel(channelName string) error {
	return inboundChannel.Remove(channelName)
}
