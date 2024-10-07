package messagesystem

/* import (
	"fmt"
	"time"

	"github.com/hex-api-go/pkg/core/infrastructure/message_system/bus"
	"github.com/hex-api-go/pkg/core/infrastructure/message_system/channel"
	"github.com/hex-api-go/pkg/core/infrastructure/message_system/channel/gochannel"
	"github.com/hex-api-go/pkg/core/infrastructure/message_system/message"
)

type MessageSystem struct {
	publisherChannels    map[string]channel.MessageOutboundChannel
	consumerChannels     map[string]channel.MessageInboundChannel
	configuredCommandBus map[string]bus.CommandBus
}

func NewMessageSystem() *MessageSystem {
	return &MessageSystem{
		publisherChannels:    map[string]channel.MessageOutboundChannel{},
		consumerChannels:     map[string]channel.MessageInboundChannel{},
		configuredCommandBus: map[string]bus.CommandBus{},
	}
}

func (m *MessageSystem) GetCommandBus() bus.CommandBus {
	return m.getCommandBus(defaultCommandChannelName)
}

func (m *MessageSystem) GetCommandBusByChannel(channelName string) bus.CommandBus {
	return m.getCommandBus(channelName)
}

func (m *MessageSystem) getCommandBus(channelName string) bus.CommandBus {

	busInstance, error := bus.GetCommandBus(channelName)
	if error != nil {
		panic(fmt.Sprintf("channel %s does not configured.", channelName))
	}
	return busInstance.(bus.CommandBus)
}

func (m *MessageSystem) Start() error {
	fmt.Println("Starting message system...")

	m.registerDefaultChannels()
	m.registerDefaultBus()
	m.startDefaultConsumers()
	return nil
}

func (m *MessageSystem) Shutdown() error {
	fmt.Println("shutting down...")
	return nil
}

func (m *MessageSystem) registerDefaultChannels() {
	fmt.Println("registering default channels...")
	gochannel.NewInternalPubSubChannelConfiguration(defaultCommandChannelName)
	fmt.Println("register default channels OK!")
}

func (m *MessageSystem) registerDefaultBus() {
	fmt.Println("registering default bus...")
	bus.NewCommandBusBuilder(defaultCommandChannelName).Build()
	fmt.Println("register default bus OK!")
}

func (m *MessageSystem) startDefaultConsumers() {
	fmt.Println("starting default consumers...")
	commandConsumer, _ := channel.GetConsumerChannel(defaultCommandChannelName)
	commandConsumer.Subscribe(func(msg any) {
		time.Sleep(time.Second * 10)
		aa := msg.(*message.GenericMessage)
		if aa.GetHeaders().GetReplyChannel() != "" {
			fmt.Println("reply channel", aa.GetHeaders().GetReplyChannel())
			replyChannel, _ := channel.GetPublisherChannel(aa.GetHeaders().GetReplyChannel())
			aa.GetHeaders().SetReplyChannel("")
			replyChannel.Send(aa)
		}
		fmt.Println("------------------------------------------------------------------")
	})
} */
