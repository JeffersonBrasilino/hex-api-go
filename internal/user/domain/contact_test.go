package domain_test

import (
	"testing"

	domain "github.com/jeffersonbrasilino/hex-api-go/internal/user/domain"
)

func TestNewContact(t *testing.T) {
	t.Run("Should success when create contact with valid data", func(t *testing.T) {
		t.Parallel()
		contactType, _ := domain.NewContactType(&domain.ContactTypeProps{
			UuId:        "1",
			Description: "1",
		})

		props := &domain.ContactProps{
			UuId:        "1",
			Description: "1",
			ContactType: contactType,
		}

		contact, err := domain.NewContact(props)
		if err != nil {
			t.Errorf("Should return a contact, got: %v", err)
		}

		if contact == nil {
			t.Error("Should return a contact, got nil")
		}
	})

	t.Run("Should fail when create contact with invalid data", func(t *testing.T) {
		t.Parallel()
		props := &domain.ContactProps{
			UuId:        "",
			Description: "",
			ContactType: nil,
		}

		contact, err := domain.NewContact(props)
		if err == nil {
			t.Errorf("Should return an error, got: %v", err)
		}

		if contact != nil {
			t.Error("Should return an error, got contact")
		}

		if err.Error() != `{"Description":{"IsValid":false,"FailedValidators":["required"]},"UuId":{"IsValid":false,"FailedValidators":["required"]}}` {
			t.Errorf("Should return an error, got: %v", err)
		}
	})
}

func TestContactGetProps(t *testing.T) {
	contactType, _ := domain.NewContactType(&domain.ContactTypeProps{
		UuId:        "1",
		Description: "1",
	})

	props := &domain.ContactProps{
		UuId:        "1",
		Description: "1",
		ContactType: contactType,
	}
	contact, _ := domain.NewContact(props)
	var cases = []struct {
		description string
		getFunc     func() any
		expected    any
	}{
		{
			description: "Should return the contact description",
			getFunc:     func() any { return contact.Description() },
			expected:    props.Description,
		},
		{
			description: "Should return the contact type",
			getFunc:     func() any { return contact.ContactType() },
			expected:    props.ContactType,
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
