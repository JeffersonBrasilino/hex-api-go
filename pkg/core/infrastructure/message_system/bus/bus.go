package bus

type (
	CommandBus interface {
		WithChannelGateway(channelName string) CommandBus
		Send(route string, message []byte, properties map[string]string) error
	}
)

var (
	commandBus *MessageSystemCommandBus
)

func Build(){
	commandBus = NewCommandBus()
}

func GetCommandBus() *MessageSystemCommandBus{
	return commandBus
}