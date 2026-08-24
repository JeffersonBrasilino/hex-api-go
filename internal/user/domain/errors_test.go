// Unit tests for the custom domain errors defined in errors.go.
//
// Intent: verify that InvalidSessionError and UserBlockedError are correctly
// constructed, carry the given message, and satisfy the error interface.
package domain_test

import (
	"testing"

	"github.com/jeffersonbrasilino/hex-api-go/internal/user/domain"
)

// TestNewInvalidSessionError verifies that NewInvalidSessionError builds an
// *InvalidSessionError whose Error() returns the provided message and that
// the returned value satisfies the error interface.
func TestNewInvalidSessionError(t *testing.T) {
	t.Run("should create error with given message", func(t *testing.T) {
		t.Parallel()
		message := "session is invalid or expired"
		err := domain.NewInvalidSessionError(message)
		if err == nil {
			t.Fatal("NewInvalidSessionError should not return nil")
		}
		if err.Error() != message {
			t.Errorf("Error() = %q, want %q", err.Error(), message)
		}
	})

	t.Run("should create error with empty message", func(t *testing.T) {
		t.Parallel()
		err := domain.NewInvalidSessionError("")
		if err == nil {
			t.Fatal("NewInvalidSessionError should not return nil")
		}
		if err.Error() != "" {
			t.Errorf("Error() = %q, want empty string", err.Error())
		}
	})

	t.Run("should implement error interface", func(t *testing.T) {
		t.Parallel()
		var _ error = domain.NewInvalidSessionError("session error")
	})
}

// TestNewUserBlockedError verifies that NewUserBlockedError builds a
// *UserBlockedError whose Error() returns the provided message and that the
// returned value satisfies the error interface.
func TestNewUserBlockedError(t *testing.T) {
	t.Run("should create error with given message", func(t *testing.T) {
		t.Parallel()
		message := "user is temporarily blocked"
		err := domain.NewUserBlockedError(message)
		if err == nil {
			t.Fatal("NewUserBlockedError should not return nil")
		}
		if err.Error() != message {
			t.Errorf("Error() = %q, want %q", err.Error(), message)
		}
	})

	t.Run("should create error with empty message", func(t *testing.T) {
		t.Parallel()
		err := domain.NewUserBlockedError("")
		if err == nil {
			t.Fatal("NewUserBlockedError should not return nil")
		}
		if err.Error() != "" {
			t.Errorf("Error() = %q, want empty string", err.Error())
		}
	})

	t.Run("should implement error interface", func(t *testing.T) {
		t.Parallel()
		var _ error = domain.NewUserBlockedError("user blocked error")
	})
}
