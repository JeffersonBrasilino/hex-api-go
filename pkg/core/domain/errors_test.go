package domain_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/hex-api-go/pkg/core/domain"
)

func TestNewDependencyError(t *testing.T) {
	t.Run("returns non-nil error with given message", func(t *testing.T) {
		t.Parallel()
		msg := "dependency unavailable"
		err := domain.NewDependencyError(msg)
		if err == nil {
			t.Fatal("NewDependencyError should not return nil")
		}
		if err.Error() != msg {
			t.Errorf("Error() should return %q, got %q", msg, err.Error())
		}
	})

	t.Run("returns error with empty message when empty string is passed", func(t *testing.T) {
		t.Parallel()
		err := domain.NewDependencyError("")
		if err == nil {
			t.Fatal("NewDependencyError should not return nil")
		}
		if err.Error() != "" {
			t.Errorf("Error() should return empty string, got %q", err.Error())
		}
	})
}

func TestNewNotFoundError(t *testing.T) {
	t.Run("returns non-nil error with given message", func(t *testing.T) {
		t.Parallel()
		msg := "resource not found"
		err := domain.NewNotFoundError(msg)
		if err == nil {
			t.Fatal("NewNotFoundError should not return nil")
		}
		if err.Error() != msg {
			t.Errorf("Error() should return %q, got %q", msg, err.Error())
		}
	})

	t.Run("returns error with empty message when empty string is passed", func(t *testing.T) {
		t.Parallel()
		err := domain.NewNotFoundError("")
		if err == nil {
			t.Fatal("NewNotFoundError should not return nil")
		}
		if err.Error() != "" {
			t.Errorf("Error() should return empty string, got %q", err.Error())
		}
	})
}

func TestNewInternalError(t *testing.T) {
	t.Run("returns non-nil error with given message", func(t *testing.T) {
		t.Parallel()
		msg := "internal server error"
		err := domain.NewInternalError(msg)
		if err == nil {
			t.Fatal("NewInternalError should not return nil")
		}
		if err.Error() != msg {
			t.Errorf("Error() should return %q, got %q", msg, err.Error())
		}
	})

	t.Run("returns error with empty message when empty string is passed", func(t *testing.T) {
		t.Parallel()
		err := domain.NewInternalError("")
		if err == nil {
			t.Fatal("NewInternalError should not return nil")
		}
		if err.Error() != "" {
			t.Errorf("Error() should return empty string, got %q", err.Error())
		}
	})
}

func TestNewAlreadyExistsError(t *testing.T) {
	t.Run("returns non-nil error with given message", func(t *testing.T) {
		t.Parallel()
		msg := "resource already exists"
		err := domain.NewAlreadyExistsError(msg)
		if err == nil {
			t.Fatal("NewAlreadyExistsError should not return nil")
		}
		if err.Error() != msg {
			t.Errorf("Error() should return %q, got %q", msg, err.Error())
		}
	})

	t.Run("returns error with empty message when empty string is passed", func(t *testing.T) {
		t.Parallel()
		err := domain.NewAlreadyExistsError("")
		if err == nil {
			t.Fatal("NewAlreadyExistsError should not return nil")
		}
		if err.Error() != "" {
			t.Errorf("Error() should return empty string, got %q", err.Error())
		}
	})
}

func TestNewInvalidDataError(t *testing.T) {
	t.Run("returns non-nil error with given message", func(t *testing.T) {
		t.Parallel()
		msg := "invalid data"
		err := domain.NewInvalidDataError(msg)
		if err == nil {
			t.Fatal("NewInvalidDataError should not return nil")
		}
		if err.Error() != msg {
			t.Errorf("Error() should return %q, got %q", msg, err.Error())
		}
	})

	t.Run("returns error with empty message when empty string is passed", func(t *testing.T) {
		t.Parallel()
		err := domain.NewInvalidDataError("")
		if err == nil {
			t.Fatal("NewInvalidDataError should not return nil")
		}
		if err.Error() != "" {
			t.Errorf("Error() should return empty string, got %q", err.Error())
		}
	})
}

func TestDomainError_Error(t *testing.T) {
	t.Run("returns only message when no previous error", func(t *testing.T) {
		t.Parallel()
		msg := "not found"
		err := domain.NewNotFoundError(msg)
		got := err.Error()
		if got != msg {
			t.Errorf("Error() should return %q, got %q", msg, got)
		}
	})

	t.Run("returns message and previous error when SetPreviousError was called", func(t *testing.T) {
		t.Parallel()
		msg := "wrap"
		prev := errors.New("cause")
		err := domain.NewNotFoundError(msg)
		err.SetPreviousError(prev)
		got := err.Error()
		if !strings.Contains(got, msg) {
			t.Errorf("Error() should contain message %q, got %q", msg, got)
		}
		if !strings.Contains(got, "previous:") {
			t.Errorf("Error() should contain 'previous:', got %q", got)
		}
		if !strings.Contains(got, "cause") {
			t.Errorf("Error() should contain previous error message, got %q", got)
		}
	})
}

func TestDomainError_SetPreviousError(t *testing.T) {
	t.Run("sets previous error and Error includes it", func(t *testing.T) {
		t.Parallel()
		msg := "top level"
		prev := errors.New("root cause")
		err := domain.NewInternalError(msg)
		err.SetPreviousError(prev)
		got := err.Error()
		if !strings.Contains(got, msg) {
			t.Errorf("Error() should contain %q, got %q", msg, got)
		}
		if !strings.Contains(got, "root cause") {
			t.Errorf("Error() should contain previous error, got %q", got)
		}
	})

	t.Run("SetPreviousError with nil does not append previous in Error", func(t *testing.T) {
		t.Parallel()
		msg := "only message"
		err := domain.NewDependencyError(msg)
		err.SetPreviousError(nil)
		got := err.Error()
		if got != msg {
			t.Errorf("Error() should return only message when previous is nil, got %q", got)
		}
	})
}
