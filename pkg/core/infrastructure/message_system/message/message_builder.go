package message

type MessageBuilder struct {
	payload       any
	route         string
	messageType   MessageType
	schema        string
	contentType   string
	replyChannel  PublisherChannel
	customHeaders customHeaders
	correlationId string
	channelName   string
}

func NewMessageBuilder() *MessageBuilder {
	return &MessageBuilder{}
}

func NewMessageBuilderFromMessage(msg *Message) *MessageBuilder {
	return &MessageBuilder{
		payload:       msg.GetPayload(),
		route:         msg.GetHeaders().Route,
		messageType:   msg.GetHeaders().MessageType,
		schema:        msg.GetHeaders().Schema,
		contentType:   msg.GetHeaders().ContentType,
		replyChannel:  msg.GetHeaders().ReplyChannel,
		customHeaders: msg.GetHeaders().CustomHeaders,
		correlationId: msg.GetHeaders().CorrelationId,
		channelName:   msg.GetHeaders().ChannelName,
	}
}

func (b *MessageBuilder) WithPayload(payload any) *MessageBuilder {
	b.payload = payload
	return b
}

func (b *MessageBuilder) WithMessageType(typeMessage MessageType) *MessageBuilder {
	b.messageType = typeMessage
	return b
}

func (b *MessageBuilder) WithRoute(route string) *MessageBuilder {
	b.route = route
	return b
}

func (b *MessageBuilder) WithSchema(value string) *MessageBuilder {
	b.schema = value
	return b
}

func (b *MessageBuilder) WithContentType(value string) *MessageBuilder {
	b.contentType = value
	return b
}

func (b *MessageBuilder) WithReplyChannel(value PublisherChannel) *MessageBuilder {
	b.replyChannel = value
	return b
}

func (b *MessageBuilder) WithCustomHeader(value customHeaders) *MessageBuilder {
	b.customHeaders = value
	return b
}

func (b *MessageBuilder) WithCorrelationId(value string) *MessageBuilder {
	b.correlationId = value
	return b
}

func (b *MessageBuilder) WithChannelName(value string) *MessageBuilder {
	b.channelName = value
	return b
}

func (b *MessageBuilder) Build() *Message {
	headers := b.buildHeaders()
	return NewMessage(b.payload, headers)
}

func (b *MessageBuilder) buildHeaders() *messageHeaders {
	headers := NewMessageHeaders(
		b.route,
		b.messageType,
		b.schema,
		b.contentType,
		b.replyChannel,
		b.correlationId,
		b.channelName,
	)
	if b.customHeaders != nil {
		headers.SetCustomHeaders(b.customHeaders)
	}
	return headers
}
