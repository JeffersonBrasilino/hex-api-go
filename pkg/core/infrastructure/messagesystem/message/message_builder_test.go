package message_test

import (
	"context"
	"testing"

	"github.com/hex-api-go/pkg/core/infrastructure/messagesystem/message"
)

func TestNewMessageBuilder(t *testing.T) {
	t.Parallel()
	b := message.NewMessageBuilder()
	if b == nil {
		t.Error("NewMessageBuilder() should not return nil")
	}
}

func TestNewMessageBuilderFromMessage(t *testing.T) {
	t.Parallel()
	b := message.NewMessageBuilder().
		WithPayload("payload").
		WithRoute("route").
		WithMessageType(1).
		WithCorrelationId("cid").
		WithChannelName("ch").
		WithReplyChannelName("rch").
		WithContext(context.Background())
	msg := b.Build()
	b2 := message.NewMessageBuilderFromMessage(msg)
	if b2 == nil {
		t.Error("NewMessageBuilderFromMessage() should not return nil")
	}
}

func TestWithPayload(t *testing.T) {
	t.Parallel()
	data := "payload"
	b := message.NewMessageBuilder().WithPayload(data).Build()
	if b.GetPayload() != data {
		t.Error("WithPayload did not set payload correctly")
	}
}

func TestWithMessageType(t *testing.T) {
	t.Parallel()
	data := message.Command
	b := message.NewMessageBuilder().WithMessageType(data).Build()
	if b.GetHeaders().MessageType != data {
		t.Error("WithMessageType did not set messageType correctly")
	}
}

func TestWithRoute(t *testing.T) {
	t.Parallel()
	b := message.NewMessageBuilder().WithRoute("route").Build()
	if b.GetHeaders().Route != "route" {
		t.Error("WithRoute did not set route correctly")
	}
}

func TestWithReplyChannel(t *testing.T) {
	t.Parallel()
	var ch message.PublisherChannel
	b := message.NewMessageBuilder().WithReplyChannel(ch).Build()
	if b.GetHeaders().ReplyChannel != ch {
		t.Error("WithReplyChannel did not set replyChannel correctly")
	}
}

func TestWithCustomHeader(t *testing.T) {
	t.Parallel()
	h := message.CustomHeaders{"k": "v"}
	b := message.NewMessageBuilder().WithCustomHeader(h).Build()
	if b.GetHeaders().CustomHeaders["k"] != "v" {
		t.Error("WithCustomHeader did not set customHeaders correctly")
	}
}

func TestWithCorrelationId(t *testing.T) {
	t.Parallel()
	b := message.NewMessageBuilder().WithCorrelationId("cid")
	if b.Build().GetHeaders().CorrelationId != "cid" {
		t.Error("WithCorrelationId did not set correlationId correctly")
	}
}

func TestWithChannelName(t *testing.T) {
	t.Parallel()
	b := message.NewMessageBuilder().WithChannelName("ch")
	if b.Build().GetHeaders().ChannelName != "ch" {
		t.Error("WithChannelName did not set channelName correctly")
	}
}

func TestWithReplyChannelName(t *testing.T) {
	t.Parallel()
	b := message.NewMessageBuilder().WithReplyChannelName("rch")
	if b.Build().GetHeaders().ReplyChannelName != "rch" {
		t.Error("WithReplyChannelName did not set replyChannelName correctly")
	}
}

func TestWithContext(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	b := message.NewMessageBuilder().WithContext(ctx).Build()
	if b.GetContext() != ctx {
		t.Error("WithContext did not set context correctly")
	}
}

func TestBuild(t *testing.T) {

	t.Parallel()
	b := message.NewMessageBuilder().
		WithPayload("payload").
		WithRoute("route").
		WithMessageType(1).
		WithCorrelationId("cid").
		WithChannelName("ch").
		WithReplyChannelName("rch").
		WithContext(context.Background())
	msg := b.Build()
	if msg == nil {
		t.Error("Build() should not return nil")
	}
}
