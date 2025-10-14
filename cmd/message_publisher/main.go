package main

import (
	"context"

	"github.com/hex-api-go/pkg/core/infrastructure/messagesystem"
	kafka "github.com/hex-api-go/pkg/core/infrastructure/messagesystem/channel/kafka"
)

type Command struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func CreateCommand(Username, Password string) *Command {
	return &Command{
		Username,
		Password,
	}
}

func (c *Command) Name() string {
	return "createUser"
}

func main() {

	//create kafka connection
	//The connection can, and is even recommended,
	//be registered only once after registration.
	//To use it in your channels, simply use its name in the channel reference name.
	messagesystem.AddChannelConnection(
		kafka.NewConnection("defaultConKafka", []string{"localhost:9093"}),
	)

	//create publisher channel
	publisherChannel := kafka.NewPublisherChannelAdapterBuilder(
		"defaultConKafka",
		"messagesystem.topic",
	)

	//register publisher channel on message system
	messagesystem.AddPublisherChannel(publisherChannel)

	//start the message system
	messagesystem.Start()

	// //get command bus
	// commandBus := messagesystem.CommandBusByChannel("messagesystem.topic")
	// commandBus.SendAsync(context.Background(), CreateCommand("teste", "123"))
	// commandBus.SendRawAsync(context.Background(), "SendAsyncRoute", "SendRawAsync command custom payload", map[string]string{"typeAction": "command"})

	//get Query bus
	queryBus := messagesystem.QueryBusByChannel("messagesystem.topic")
	queryBus.SendAsync(context.Background(), CreateCommand("teste", "123"))
	queryBus.SendRawAsync(context.Background(), "SendAsyncRoute", "SendRawAsync query custom payload", map[string]string{"typeAction": "query"})

	// //get event bus
	// eventBus := messagesystem.EventBusByChannel("messagesystem.topic")
	// eventBus.Publish(context.Background(), CreateCommand("teste", "123"))
	// eventBus.PublishRaw(context.Background(), "SendAsyncRoute", "SendRawAsync query custom payload", map[string]string{"aaaa": "uheuehuehueh"})

	//message system graceful shutdown
	messagesystem.Shutdown()
}
