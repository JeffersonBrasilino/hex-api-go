package bus

/* var createdEventBus sync.Map

type EventBus struct {
	*messageBus
}

func NewEventBus(gateway *endpoint.Gateway, channelName string) *EventBus {

	bus, ok := createdEventBus.Load(channelName)
	if ok {
		return bus.(*EventBus)
	}

	eventBus := &EventBus{
		messageBus: &messageBus{
			gateway,
		},
	}
	createdEventBus.Store(channelName, bus)
	return eventBus
}

func (c *EventBus) Publish(ctx context.Context, action handler.Action) error {
	builder := c.buildMessage()
	msg := builder.WithPayload(action).
		WithRoute(action.Name()).
		Build()
	return c.publishMessage(ctx, msg)
}

func (c *EventBus) PublishRaw(ctx context.Context, route string, payload any, headers map[string]string) error {
	builder := c.buildMessage()
	msg := builder.WithPayload(payload).
		WithRoute(route).
		WithCustomHeader(headers).
		Build()
	return c.publishMessage(ctx, msg)
}

func (c *EventBus) buildMessage() *message.MessageBuilder {
	builder := message.NewMessageBuilder().
		WithMessageType(message.Event).
		WithCorrelationId(uuid.New().String())
	return builder
}
*/