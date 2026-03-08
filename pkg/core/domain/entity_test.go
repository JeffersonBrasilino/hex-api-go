package domain_test

import (
	"testing"

	"github.com/hex-api-go/pkg/core/domain"
)

func TestNewEntity(t *testing.T) {
	t.Run("returns non-nil entity with given uuid", func(t *testing.T) {
		t.Parallel()
		uuid := "550e8400-e29b-41d4-a716-446655440000"
		e := domain.NewEntity(uuid)
		if e == nil {
			t.Fatal("NewEntity should not return nil")
		}
		if e.Uuid() != uuid {
			t.Errorf("entity uuid should be %q, got %q", uuid, e.Uuid())
		}
	})

	t.Run("returns entity with empty uuid when empty string is passed", func(t *testing.T) {
		t.Parallel()
		e := domain.NewEntity("")
		if e == nil {
			t.Fatal("NewEntity should not return nil")
		}
		if e.Uuid() != "" {
			t.Errorf("entity uuid should be empty, got %q", e.Uuid())
		}
	})
}

func TestEntity_Uuid(t *testing.T) {
	t.Run("returns uuid set at creation", func(t *testing.T) {
		t.Parallel()
		uuid := "abc-123"
		e := domain.NewEntity(uuid)
		got := e.Uuid()
		if got != uuid {
			t.Errorf("Uuid() should return %q, got %q", uuid, got)
		}
	})

	t.Run("returns empty string when entity was created with empty uuid", func(t *testing.T) {
		t.Parallel()
		e := domain.NewEntity("")
		got := e.Uuid()
		if got != "" {
			t.Errorf("Uuid() should return empty string, got %q", got)
		}
	})
}
