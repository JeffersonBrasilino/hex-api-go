package domain_test

import (
	"testing"

	"github.com/jeffersonbrasilino/hex-api-go/internal/user/domain"
)

func TestNewContactType(t *testing.T) {
	t.Run("should create contact type successfully", func(t *testing.T) {
		t.Parallel()

		props := &domain.ContactTypeProps{
			UuId:        "123e4567-e89b-12d3-a456-426614174000",
			Description: "Email",
		}

		contactType, err := domain.NewContactType(props)
		if err != nil {
			t.Errorf("expected no error, got: %v", err)
		}
		if contactType == nil {
			t.Fatal("expected contact type to be created")
		}
		if contactType.Description() != "Email" {
			t.Errorf("expected description to be 'Email', got: %s", contactType.Description())
		}
	})

	t.Run("should return error when uuid is missing", func(t *testing.T) {
		t.Parallel()

		props := &domain.ContactTypeProps{
			UuId:        "", // Invalid
			Description: "Email",
		}

		contactType, err := domain.NewContactType(props)
		if err == nil {
			t.Error("expected validation error, got nil")
		}
		if contactType != nil {
			t.Error("expected contact type to be nil on error")
		}
	})

}

func TestContactType_Description(t *testing.T) {
	t.Run("should return correct description", func(t *testing.T) {
		t.Parallel()

		props := &domain.ContactTypeProps{
			UuId:        "123e4567-e89b-12d3-a456-426614174000",
			Description: "Phone",
		}

		contactType, err := domain.NewContactType(props)
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		desc := contactType.Description()
		if desc != "Phone" {
			t.Errorf("expected 'Phone', got: %s", desc)
		}
	})
}
