package message

import (
	"context"
	"testing"
)

func TestMessageType_String(t *testing.T) {
	if Command.String() != "Command" {
		t.Error("Command.String() should return 'Command'")
	}
	if Query.String() != "Query" {
		t.Error("Query.String() should return 'Query'")
	}
	if Event.String() != "Event" {
		t.Error("Event.String() should return 'Event'")
	}
	if Document.String() != "Document" {
		t.Error("Document.String() should return 'Document'")
	}
	var mt MessageType = 99
	if mt.String() != "Message" {
		t.Error("Unknown MessageType should return 'Message'")
	}
}

func TestNewMessageHeaders_Success(t *testing.T) {
	headers := NewMessageHeaders("route", Command, nil, "cid", "ch", "rch")
	if headers.Route != "route" {
		t.Error("Route not set correctly")
	}
	if headers.MessageType != Command {
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

func TestSetCustomHeaders_Success(t *testing.T) {
	headers := NewMessageHeaders("route", Command, nil, "cid", "ch", "rch")
	ch := customHeaders{"foo": "bar"}
	headers.SetCustomHeaders(ch)
	if headers.CustomHeaders["foo"] != "bar" {
		t.Error("SetCustomHeaders did not set value correctly")
	}
}

func TestMessageHeaders_MarshalJSON_Success(t *testing.T) {
	headers := NewMessageHeaders("route", Command, nil, "cid", "ch", "rch")
	headers.SetCustomHeaders(customHeaders{"foo": "bar"})
	_, err := headers.MarshalJSON()
	if err != nil {
		t.Error("MarshalJSON should not return error")
	}
}

func TestNewMessageAndGetters_Success(t *testing.T) {
	headers := NewMessageHeaders("route", Command, nil, "cid", "ch", "rch")
	ctx := context.Background()
	msg := NewMessage("payload", headers, ctx)
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

func TestMessage_SetContext_Success(t *testing.T) {
	headers := NewMessageHeaders("route", Command, nil, "cid", "ch", "rch")
	msg := NewMessage("payload", headers, nil)
	ctx := context.Background()
	msg.SetContext(ctx)
	if msg.GetContext() != ctx {
		t.Error("SetContext did not set context correctly")
	}
}

func TestMessage_ReplyRequired(t *testing.T) {
	headers := NewMessageHeaders("route", Command, nil, "cid", "ch", "rch")
	msg := NewMessage("payload", headers, nil)
	if !msg.ReplyRequired() {
		t.Error("ReplyRequired should be true for Command")
	}
	headers.MessageType = Query
	if !msg.ReplyRequired() {
		t.Error("ReplyRequired should be true for Query")
	}
	headers.MessageType = Event
	if msg.ReplyRequired() {
		t.Error("ReplyRequired should be false for Event")
	}
	headers.MessageType = Document
	if msg.ReplyRequired() {
		t.Error("ReplyRequired should be false for Document")
	}
}

func TestMessage_MarshalJSON_Success(t *testing.T) {
	headers := NewMessageHeaders("route", Command, nil, "cid", "ch", "rch")
	msg := NewMessage("payload", headers, nil)
	_, err := msg.MarshalJSON()
	if err != nil {
		t.Error("MarshalJSON should not return error")
	}
}
