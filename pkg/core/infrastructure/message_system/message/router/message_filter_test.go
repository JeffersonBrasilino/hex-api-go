package router

import (
	"context"
	"testing"

	"github.com/hex-api-go/pkg/core/infrastructure/message_system/message"
)

func TestNewMessageFilter(t *testing.T) {
	f := NewMessageFilter(func(m message.Message) bool { return true })
	if f == nil {
		t.Error("NewMessageFilter deve retornar uma instância não nula")
	}
}

func TestMessageFilter_Handle_Pass(t *testing.T) {
	filter := NewMessageFilter(func(m message.Message) bool { return true })
	msg := &message.Message{}
	result, err := filter.Handle(context.Background(), msg)
	if err != nil {
		t.Errorf("Handle deve retornar erro nil, retornou: %v", err)
	}
	if result != msg {
		t.Error("Handle deve retornar a mensagem original quando o filtro retorna true")
	}
}

func TestMessageFilter_Handle_Block(t *testing.T) {
	filter := NewMessageFilter(func(m message.Message) bool { return false })
	msg := &message.Message{}
	result, err := filter.Handle(context.Background(), msg)
	if err != nil {
		t.Errorf("Handle deve retornar erro nil, retornou: %v", err)
	}
	if result != nil {
		t.Error("Handle deve retornar nil quando o filtro retorna false")
	}
}
