package domain_test

import (
	"testing"

	domain "github.com/jeffersonbrasilino/hex-api-go/internal/user/domain"
)

func TestNewDocument(t *testing.T) {
	t.Run("Should success when create document with valid data", func(t *testing.T) {
		t.Parallel()
		props := &domain.DocumentProps{
			Value: "1",
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
		Value: "1",
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
