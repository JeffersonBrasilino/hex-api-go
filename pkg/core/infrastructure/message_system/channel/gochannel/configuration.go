package gochannel

type configuration struct {
	channelName    string
	buildConsumer  bool
	buildPublisher bool
}

func NewPubSubChannelConfiguration(
	channelName string,
) *configuration {
	return &configuration{
		channelName:    channelName,
		buildPublisher: true,
		buildConsumer:  true,
	}
}

func NewPublisherConfiguration(
	channelName string,
) *configuration {
	return &configuration{
		channelName:    channelName,
		buildPublisher: true,
		buildConsumer:  false,
	}
}

func NewConsumerConfiguration(
	channelName string,
) *configuration {
	return &configuration{
		channelName:    channelName,
		buildPublisher: false,
		buildConsumer:  true,
	}
}

func (g *configuration) ChannelName() string {
	return g.channelName
}
