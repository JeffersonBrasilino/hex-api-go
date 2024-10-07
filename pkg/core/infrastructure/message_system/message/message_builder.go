package message

type MessageBuilder struct {
	payload       []byte
	route         string
	messageType   string
	schema        string
	contentType   string
	replyChannel  string
	errorChannel  string
	version       string
	customHeaders customHeaders
	correlationId string
	channelName   string
}

func NewMessageBuilder() *MessageBuilder {
	return &MessageBuilder{}
}

func NewMessageBuilderFromMessage(msg Message) *MessageBuilder {
	return &MessageBuilder{
		payload:       msg.GetPayload(),
		route:         msg.GetHeaders().route,
		messageType:   msg.GetHeaders().messageType,
		schema:        msg.GetHeaders().schema,
		contentType:   msg.GetHeaders().contentType,
		replyChannel:  msg.GetHeaders().replyChannel,
		errorChannel:  msg.GetHeaders().errorChannel,
		version:       msg.GetHeaders().version,
		customHeaders: msg.GetHeaders().customHeaders,
		correlationId: msg.GetHeaders().correlationId,
		channelName:   msg.GetHeaders().channelName,
	}
}

func (b *MessageBuilder) WithPayload(payload []byte) *MessageBuilder {
	b.payload = payload
	return b
}

func (b *MessageBuilder) WithMessageType(typeMessage MessageType) *MessageBuilder {
	b.messageType = typeMessage.String()
	return b
}

func (b *MessageBuilder) WithRoute(route string) *MessageBuilder {
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

func (b *MessageBuilder) WithReplyChannel(value string) *MessageBuilder {
	b.replyChannel = value
	return b
}

func (b *MessageBuilder) WithErrorChannel(value string) *MessageBuilder {
	b.errorChannel = value
	return b
}

func (b *MessageBuilder) WithVersion(value string) *MessageBuilder {
	b.version = value
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

func (b *MessageBuilder) Build() *GenericMessage {
	headers := b.buildHeaders()
	return NewGenericMessage(b.payload, headers)
}

func (b *MessageBuilder) buildHeaders() *messageHeaders {
	headers := NewMessageHeaders(
		b.route,
		b.messageType,
		b.schema,
		b.contentType,
		b.replyChannel,
		b.errorChannel,
		b.version,
		b.correlationId,
		b.channelName,
	)
	if b.customHeaders != nil {
		headers.SetCustomHeaders(b.customHeaders)
	}
	return headers
}
