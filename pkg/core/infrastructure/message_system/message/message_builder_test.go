package message

import (
	"context"
	"testing"
)

func TestNewMessageBuilder(t *testing.T) {
	b := NewMessageBuilder()
	if b == nil {
		t.Error("NewMessageBuilder() should not return nil")
	}
}

func TestNewMessageBuilderFromMessage(t *testing.T) {
	b := NewMessageBuilder().
		WithPayload("payload").
		WithRoute("route").
		WithMessageType(1).
		WithCorrelationId("cid").
		WithChannelName("ch").
		WithReplyChannelName("rch").
		WithContext(context.Background())
	msg := b.Build()
	b2 := NewMessageBuilderFromMessage(msg)
	if b2 == nil {
		t.Error("NewMessageBuilderFromMessage() should not return nil")
	}
}

func TestWithPayload(t *testing.T) {
	b := NewMessageBuilder().WithPayload("payload")
	if b.payload != "payload" {
		t.Error("WithPayload did not set payload correctly")
	}
}

func TestWithMessageType(t *testing.T) {
	b := NewMessageBuilder().WithMessageType(2)
	if b.messageType != 2 {
		t.Error("WithMessageType did not set messageType correctly")
	}
}

func TestWithRoute(t *testing.T) {
	b := NewMessageBuilder().WithRoute("route")
	if b.route != "route" {
		t.Error("WithRoute did not set route correctly")
	}
}

func TestWithReplyChannel(t *testing.T) {
	var ch PublisherChannel
	b := NewMessageBuilder().WithReplyChannel(ch)
	if b.replyChannel != ch {
		t.Error("WithReplyChannel did not set replyChannel correctly")
	}
}

func TestWithCustomHeader(t *testing.T) {
	var h customHeaders = customHeaders{"k": "v"}
	b := NewMessageBuilder().WithCustomHeader(h)
	if b.customHeaders["k"] != "v" {
		t.Error("WithCustomHeader did not set customHeaders correctly")
	}
}

func TestWithCorrelationId(t *testing.T) {
	b := NewMessageBuilder().WithCorrelationId("cid")
	if b.correlationId != "cid" {
		t.Error("WithCorrelationId did not set correlationId correctly")
	}
}

func TestWithChannelName(t *testing.T) {
	b := NewMessageBuilder().WithChannelName("ch")
	if b.channelName != "ch" {
		t.Error("WithChannelName did not set channelName correctly")
	}
}

func TestWithReplyChannelName(t *testing.T) {
	b := NewMessageBuilder().WithReplyChannelName("rch")
	if b.replyChannelName != "rch" {
		t.Error("WithReplyChannelName did not set replyChannelName correctly")
	}
}

func TestWithContext(t *testing.T) {
	ctx := context.Background()
	b := NewMessageBuilder().WithContext(ctx)
	if b.context != ctx {
		t.Error("WithContext did not set context correctly")
	}
}

func TestBuild(t *testing.T) {
	b := NewMessageBuilder().
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

func TestBuildHeaders(t *testing.T) {
	b := NewMessageBuilder().
		WithRoute("route").
		WithMessageType(1).
		WithCorrelationId("cid").
		WithChannelName("ch").
		WithReplyChannelName("rch")
	headers := b.buildHeaders()
	if headers == nil {
		t.Error("buildHeaders() should not return nil")
	}
}
