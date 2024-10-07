package messagingo

import (
	"github.com/hex-api-go/pkg/messagingo/channel"
	"github.com/hex-api-go/pkg/messagingo/client"
	"github.com/hex-api-go/pkg/messagingo/container"
	"github.com/hex-api-go/pkg/messagingo/handler"
)

type MessageSystem struct {
	container                      container.Container
	channelPublisherConfigurations []channel.ChannelPublisherConfiguration
	messageConverter               MessageConverter
	afterSendMiddleware            []handler.MessageHandler
	beforeSendMiddleware           []handler.MessageHandler
}

const (
	connectionClientContainerKey = "defaultClientConnection"
)

func CreateMessageSystem() *MessageSystem {
	return &MessageSystem{
		container:                      container.CreateInMemoryContainer(),
		channelPublisherConfigurations: []channel.ChannelPublisherConfiguration{},
		afterSendMiddleware:            []handler.MessageHandler{},
		beforeSendMiddleware:           []handler.MessageHandler{},
	}
}

func (s *MessageSystem) WithConnectionClient(
	connectionClient client.ConnectionClient,
) *MessageSystem {
	if s.container.Has(connectionClientContainerKey) {
		panic("already exists connection to broker")
	}
	s.container.Set(connectionClientContainerKey, connectionClient)
	return s
}

func (m *MessageSystem) WithMessageConverter(converter MessageConverter) *MessageSystem {
	m.messageConverter = converter
	return m
}

func (m *MessageSystem) WithBeforeSendMiddleware(middleware handler.MessageHandler) *MessageSystem {
	m.beforeSendMiddleware = append(m.beforeSendMiddleware, middleware)
	return m
}

func (m *MessageSystem) WithAfterSendMiddleware(middleware handler.MessageHandler) *MessageSystem {
	m.afterSendMiddleware = append(m.afterSendMiddleware, middleware)
	return m
}

/* func (m *MessageSystem) RegisterPublisherChannel(
	channelConfiguration channel.ChannelPublisherConfiguration,
) *MessageSystem {
	if m.container.Has(channelConfiguration.GetQueueOrTopicName()) {
		panic(
			fmt.Sprintf(
				"There is already a publisher for %s",
				channelConfiguration.GetQueueOrTopicName()),
		)
	}
	connectionInstance := m.container.Get(connectionClientContainerKey).(client.ConnectionClient)
	buildedChannel := channelConfiguration.Build(connectionInstance)
	m.container.Set(channelConfiguration.GetQueueOrTopicName(), buildedChannel)
	return m
}

func (m *MessageSystem) buildGateways() {
	//handlerPipeline := m.buildMessageHandlerPipeline(nil)

} */

func (m *MessageSystem) Start() *ConfiguredMessageSystem {
	return &ConfiguredMessageSystem{m.container}
}
