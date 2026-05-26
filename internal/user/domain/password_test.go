package domain_test

import (
	"testing"

	"github.com/jeffersonbrasilino/hex-api-go/internal/user/domain"
)

func TestNewPassword(t *testing.T) {
	t.Run("should succeed when creating password with valid strong value", func(t *testing.T) {
		t.Parallel()
		props := &domain.PasswordProps{
			Value: "StrongP@ss123",
		}

		pwd, err := domain.NewPassword(props)
		if err != nil {
			t.Errorf("NewPassword should return nil error, got: %v", err)
		}

		if pwd == nil {
			t.Error("NewPassword should return a valid pointer to Password, got nil")
		}
	})

	t.Run("should fail when value is empty (required)", func(t *testing.T) {
		t.Parallel()
		props := &domain.PasswordProps{
			Value: "",
		}

		pwd, err := domain.NewPassword(props)
		if err == nil {
			t.Error("NewPassword should return an error when Value is empty")
		}
		if pwd != nil {
			t.Error("NewPassword should return nil Password on error")
		}

		expectedErr := `{"Value":{"IsValid":false,"FailedValidators":["required"]}}`
		if err.Error() != expectedErr {
			t.Errorf("Expected error %s, got: %v", expectedErr, err.Error())
		}
	})

	t.Run("should fail when missing lowercase letter", func(t *testing.T) {
		t.Parallel()
		props := &domain.PasswordProps{
			Value: "STRONGP@SS123",
		}

		pwd, err := domain.NewPassword(props)
		if err == nil {
			t.Error("NewPassword should return an error when missing lowercase")
		}
		if pwd != nil {
			t.Error("NewPassword should return nil Password on error")
		}

		expectedErr := `{"Value":{"FailedValidators":["strong_password"],"IsValid":false}}`
		if err.Error() != expectedErr {
			t.Errorf("Expected error %s, got: %v", expectedErr, err.Error())
		}
	})

	t.Run("should fail when missing uppercase letter", func(t *testing.T) {
		t.Parallel()
		props := &domain.PasswordProps{
			Value: "strongp@ss123",
		}

		pwd, err := domain.NewPassword(props)
		if err == nil {
			t.Error("NewPassword should return an error when missing uppercase")
		}
		if pwd != nil {
			t.Error("NewPassword should return nil Password on error")
		}
	})

	t.Run("should fail when missing digit", func(t *testing.T) {
		t.Parallel()
		props := &domain.PasswordProps{
			Value: "StrongP@ss",
		}

		pwd, err := domain.NewPassword(props)
		if err == nil {
			t.Error("NewPassword should return an error when missing digit")
		}
		if pwd != nil {
			t.Error("NewPassword should return nil Password on error")
		}
	})

	t.Run("should fail when missing special character", func(t *testing.T) {
		t.Parallel()
		props := &domain.PasswordProps{
			Value: "StrongPass123",
		}

		pwd, err := domain.NewPassword(props)
		if err == nil {
			t.Error("NewPassword should return an error when missing special character")
		}
		if pwd != nil {
			t.Error("NewPassword should return nil Password on error")
		}
	})

	t.Run("should fail when password is less than 8 characters", func(t *testing.T) {
		t.Parallel()
		props := &domain.PasswordProps{
			Value: "St@123",
		}

		pwd, err := domain.NewPassword(props)
		if err == nil {
			t.Error("NewPassword should return an error when password is too short")
		}
		if pwd != nil {
			t.Error("NewPassword should return nil Password on error")
		}
	})
}

func TestPassword_Value(t *testing.T) {
	t.Run("should return the correct password value string", func(t *testing.T) {
		t.Parallel()
		expectedValue := "StrongP@ss123"
		props := &domain.PasswordProps{
			Value: expectedValue,
		}

		pwd, _ := domain.NewPassword(props)

		if pwd.Value() != expectedValue {
			t.Errorf("Value() should return '%s', got: '%s'", expectedValue, pwd.Value())
		}
	})
}
