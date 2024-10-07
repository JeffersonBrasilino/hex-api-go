package bus

type CommandBusBuilder struct {
	topicOrQueueName string
}

func NewCommandBusBuilder(topicOrQueueName string) *CommandBusBuilder {
	return &CommandBusBuilder{
		topicOrQueueName: topicOrQueueName,
	}
}

func (b *CommandBusBuilder) Build() {
/* 	fmt.Println(fmt.Printf("building bus for %s ", b.topicOrQueueName))
	buildedGateway := gateway.NewGatewayBuilder(
		b.topicOrQueueName,
		b.topicOrQueueName,
	).WithMessageConverterFor(message.Command).Build()
	commandBus := NewCommandBus(b.topicOrQueueName, *buildedGateway)
	fmt.Println(fmt.Printf("building bus for %s OK! ", b.topicOrQueueName))
	RegisterCommandBus(commandBus) */
}
