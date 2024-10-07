package messagesystem

import (
	"fmt"
	"time"

	"github.com/hex-api-go/pkg/core/infrastructure/message_system/bus"
	"github.com/hex-api-go/pkg/core/infrastructure/message_system/channel"
	"github.com/hex-api-go/pkg/core/infrastructure/message_system/channel/gochannel"
	"github.com/hex-api-go/pkg/core/infrastructure/message_system/gateway"
	"github.com/hex-api-go/pkg/core/infrastructure/message_system/message"
)

const (
	defaultCommandChannelName = "messagesystem.command"
)

var modules = map[string]func(){
	"channel": channel.Build,
	"gateway": gateway.Build,
	"bus":     bus.Build,
}

func AddModule(name string, module func()) {
	modules[name] = module
}

func buildMessageSystemModules() {
	fmt.Println("load modules...")
	for _, build := range modules {
		build()
	}
}

func registerDefualtChannels() {
	gochannel.RegisterChannel(
		gochannel.NewPubSubChannelConfiguration(defaultCommandChannelName),
	)
}

func Start() {
	fmt.Println("starting message system...")
	registerDefualtChannels()
	buildMessageSystemModules()
	startDefaultConsumers()
}

func Shutdown() {
	fmt.Println("shutting down...")
}

func GetCommandBus() bus.MessageSystemCommandBus {
	cb := *bus.GetCommandBus()
	cb.WithChannelGateway(defaultCommandChannelName)
	return cb
}

func startDefaultConsumers() {
	fmt.Println("starting default consumers...")
	commandConsumer, _ := channel.GetInboundChannel(defaultCommandChannelName)
	commandConsumer.Subscribe(func(msg any) {
		time.Sleep(time.Second * 10)
		aa := msg.(*message.GenericMessage)
		if aa.GetHeaders().GetReplyChannel() != "" {
			fmt.Println("reply channel", aa.GetHeaders().GetReplyChannel())
			replyChannel, _ := channel.GetOutboundChannel(aa.GetHeaders().GetReplyChannel())
			aa.GetHeaders().SetReplyChannel("")
			replyChannel.Send(aa)
		}
		fmt.Println("------------------------------------------------------------------")
	})
}