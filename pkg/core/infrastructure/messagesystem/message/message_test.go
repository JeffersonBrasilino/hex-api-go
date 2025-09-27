package message_test

import (
	"context"
	"testing"

	"github.com/hex-api-go/pkg/core/infrastructure/messagesystem/message"
)

func TestMessageTypeString(t *testing.T) {
	cases := []struct {
		description string
		should      message.MessageType
		want        string
	}{
		{"should Command", message.Command, "Command"},
		{"should Query", message.Query, "Query"},
		{"should Event", message.Event, "Event"},
		{"should Document", message.Document, "Document"},
		{"should Unknown", message.MessageType(99), "Message"},
	}

	for _, c := range cases {
		t.Run(c.description, func(t *testing.T) {
			t.Parallel()
			if c.should.String() != c.want {
				t.Errorf("%s.String() should return '%s'", c.should, c.want)
			}
		})
	}
}

func TestNewMessageHeaders(t *testing.T) {
	t.Parallel()
	headers := message.NewMessageHeaders("route", message.Command, nil, "cid", "ch", "rch")
	if headers.Route != "route" {
		t.Error("Route not set correctly")
	}
	if headers.MessageType != message.Command {
		t.Error("MessageType not set correctly")
	}
	if headers.CorrelationId != "cid" {
		t.Error("CorrelationId not set correctly")
	}
	if headers.ChannelName != "ch" {
		t.Error("ChannelName not set correctly")
	}
	if headers.ReplyChannelName != "rch" {
		t.Error("ReplyChannelName not set correctly")
	}
	if headers.CustomHeaders == nil {
		t.Error("CustomHeaders should be initialized")
	}
}

func TestMessageHeaders_SetCustomHeaders(t *testing.T) {
	t.Parallel()
	headers := message.NewMessageHeaders("route", message.Command, nil, "cid", "ch", "rch")
	headers.SetCustomHeaders(message.CustomHeaders{"foo": "bar"})
	if headers.CustomHeaders["foo"] != "bar" {
		t.Error("CustomHeaders not set correctly")
	}
}

func TestCustomMessageHeaders_MarshalJSON(t *testing.T) {
	t.Parallel()
	headers := message.NewMessageHeaders("route", message.Command, nil, "cid", "ch", "rch")
	headers.SetCustomHeaders(message.CustomHeaders{"foo": "bar"})
	_, err := headers.MarshalJSON()
	if err != nil {
		t.Error("MarshalJSON should not return error")
	}
}

func TestNewMessage(t *testing.T) {
	t.Parallel()
	headers := message.NewMessageHeaders("route", message.Command, nil, "cid", "ch", "rch")
	ctx := context.Background()
	msg := message.NewMessage("payload", headers, ctx)
	if msg.GetPayload() != "payload" {
		t.Error("GetPayload did not return correct value")
	}
	if msg.GetHeaders() != headers {
		t.Error("GetHeaders did not return correct value")
	}
	if msg.GetContext() != ctx {
		t.Error("GetContext did not return correct value")
	}
}

func TestMessage_SetContext(t *testing.T) {
	t.Parallel()
	headers := message.NewMessageHeaders("route", message.Command, nil, "cid", "ch", "rch")
	msg := message.NewMessage("payload", headers, nil)
	ctx := context.Background()
	msg.SetContext(ctx)
	if msg.GetContext() != ctx {
		t.Error("SetContext did not set context correctly")
	}
}

func TestMessage_Getters(t *testing.T) {
	headers := message.NewMessageHeaders("route", message.Command, nil, "cid", "ch", "rch")
	ctx := context.Background()
	msg := message.NewMessage("payload", headers, ctx)
	cases := []struct {
		description string
		should      func() any
		want        any
	}{
		{"GetPayload should return correct value", func() any {
			return msg.GetPayload()
		}, "payload"},
		{"GetHeaders should return correct value", func() any {
			return msg.GetHeaders()
		}, headers},
		{"GetContext should return correct value", func() any {
			return msg.GetContext()
		}, ctx},
	}

	for _, c := range cases {
		t.Run(c.description, func(t *testing.T) {
			t.Parallel()
			if c.should() != c.want {
				t.Errorf("%s() should return '%v'", c.description, c.want)
			}
		})
	}
}

func TestMessage_ReplyRequired(t *testing.T) {
	cases := []struct {
		description string
		should      message.MessageType
		want        bool
	}{
		{"ReplyRequired should return true for Command", message.Command, true},
		{"ReplyRequired should return true for Query", message.Query, true},
		{"ReplyRequired should return false for Event", message.Event, false},
		{"ReplyRequired should return false for Document", message.Document, false},
	}
	for _, c := range cases {
		t.Run(c.description, func(t *testing.T) {
			t.Parallel()
			headers := message.NewMessageHeaders("route", c.should, nil, "cid", "ch", "rch")
			msg := message.NewMessage("payload", headers, nil)
			if msg.ReplyRequired() != c.want {
				t.Errorf("%q: got %v, want %v", c.description, msg.ReplyRequired(), c.want)
			}
		})
	}
}

func TestMessage_MarshalJSON(t *testing.T) {
	t.Parallel()
	headers := message.NewMessageHeaders("route", message.Command, nil, "cid", "ch", "rch")
	msg := message.NewMessage("payload", headers, nil)
	_, err := msg.MarshalJSON()
	if err != nil {
		t.Error("MarshalJSON should not return error")
	}
}