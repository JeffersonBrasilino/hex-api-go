package domain_test

import (
	"testing"

	domain "github.com/jeffersonbrasilino/hex-api-go/internal/user/domain"
)

func TestNewDocument(t *testing.T) {
	t.Run("Should success when create document with valid data", func(t *testing.T) {
		t.Parallel()
		props := &domain.DocumentProps{
			Value: "52998224725",
		}

		document, err := domain.NewDocument(props)
		if err != nil {
			t.Errorf("Should return a document, got: %v", err)
		}

		if document == nil {
			t.Error("Should return a document, got nil")
		}
	})

	t.Run("Should fail when create document with invalid data", func(t *testing.T) {
		t.Parallel()
		props := &domain.DocumentProps{
			Value: "",
		}

		document, err := domain.NewDocument(props)
		if err == nil {
			t.Errorf("Should return an error, got: %v", err)
		}

		if document != nil {
			t.Error("Should return an error, got document")
		}

		if err.Error() != `{"Value":{"IsValid":false,"FailedValidators":["required"]}}` {
			t.Errorf("Should return an error, got: %v", err)
		}
	})
}

func TestDocumentGetProps(t *testing.T) {
	props := &domain.DocumentProps{
		Value: "52998224725",
	}
	document, _ := domain.NewDocument(props)
	var cases = []struct {
		description string
		getFunc     func() any
		expected    any
	}{
		{
			description: "Should return the document value",
			getFunc:     func() any { return document.Value() },
			expected:    props.Value,
		},
	}

	for _, c := range cases {
		t.Run(c.description, func(t *testing.T) {
			t.Parallel()
			if c.getFunc() != c.expected {
				t.Errorf("Should return %v, got: %v", c.expected, c.getFunc())
			}
		})
	}
}

func TestIsValidCPF(t *testing.T) {
	cases := []struct {
		description string
		cpf         string
		wantErr     bool
		errMsg      string
	}{
		{
			description: "Should succeed with a valid CPF",
			cpf:         "52998224725",
			wantErr:     false,
		},
		{
			description: "Should fail when all digits are the same",
			cpf:         "11111111111",
			wantErr:     true,
			errMsg:      "Invalid cpf.",
		},
		{
			description: "Should fail when all digits are zero",
			cpf:         "00000000000",
			wantErr:     true,
			errMsg:      "Invalid cpf.",
		},
		{
			description: "Should fail when CPF has fewer than 11 digits",
			cpf:         "1234567890",
			wantErr:     true,
			errMsg:      "Invalid cpf.",
		},
		{
			description: "Should fail when CPF has more than 11 digits",
			cpf:         "123456789012",
			wantErr:     true,
			errMsg:      "Invalid cpf.",
		},
		{
			description: "Should fail when check digits are invalid",
			cpf:         "12345678901",
			wantErr:     true,
			errMsg:      "Invalid cpf.",
		},
	}

	for _, c := range cases {
		t.Run(c.description, func(t *testing.T) {
			t.Parallel()
			_, err := domain.NewDocument(&domain.DocumentProps{Value: c.cpf})
			if c.wantErr {
				if err == nil {
					t.Errorf("Expected error but got nil")
					return
				}
				if err.Error() != c.errMsg {
					t.Errorf("Expected error %q, got %q", c.errMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error but got: %v", err)
				}
			}
		})
	}
}
