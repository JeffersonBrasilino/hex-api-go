package domain_test

import (
	"fmt"
	"testing"

	domain "github.com/jeffersonbrasilino/hex-api-go/internal/user/domain"
)

func TestNewPerson(t *testing.T) {
	t.Run("Should success when create person with valid data", func(t *testing.T) {
		t.Parallel()
		ctt, _ := domain.NewContact(&domain.ContactProps{
			UuId:        "1",
			Description: "1",
			ContactType: "1",
		})
		dcm, _ := domain.NewDocument(&domain.DocumentProps{
			Value: "1",
		})
		props := &domain.PersonProps{
			UuId:      "1",
			Name:      "John Doe",
			BirthDate: "2000-01-01",
			Contacts: []*domain.Contact{
				ctt,
			},
			Document: dcm,
		}

		person, err := domain.NewPerson(props)
		if err != nil {
			t.Errorf("Should return a person, got: %v", err)
		}

		if person == nil {
			t.Error("Should return a person, got nil")
		}
	})

	t.Run("Should fail when create person with invalid data", func(t *testing.T) {
		t.Parallel()
		props := &domain.PersonProps{
			UuId:      "",
			Name:      "",
			BirthDate: "",
			Contacts:  make([]*domain.Contact, 0),
			Document:  nil,
		}

		person, err := domain.NewPerson(props)
		fmt.Println("PROPS", props)
		fmt.Println("RESULTADO", err)
		if err == nil {
			t.Errorf("Should return an error, got: %v", err)
		}

		if person != nil {
			t.Error("Should return an error, got person")
		}
	})
}

func TestPersonGetters(t *testing.T) {
	t.Run("Should return the correct person data", func(t *testing.T) {
		t.Parallel()
		ctt, _ := domain.NewContact(&domain.ContactProps{
			UuId:        "1",
			Description: "1",
			ContactType: "1",
		})
		dcm, _ := domain.NewDocument(&domain.DocumentProps{
			Value: "1",
		})
		props := &domain.PersonProps{
			UuId:      "1",
			Name:      "John Doe",
			BirthDate: "2000-01-01",
			Contacts: []*domain.Contact{
				ctt,
			},
			Document: dcm,
		}

		person, _ := domain.NewPerson(props)

		if person.Name() != "John Doe" {
			t.Errorf("Should return the correct name, got: %v", person.Name())
		}

		if person.BirthDate() != "2000-01-01" {
			t.Errorf("Should return the correct birth date, got: %v", person.BirthDate())
		}

		if person.Contacts()[0].Description() != "1" {
			t.Errorf("Should return the correct contact value, got: %v", person.Contacts()[0].Description())
		}

		if person.Document().Value() != "1" {
			t.Errorf("Should return the correct document value, got: %v", person.Document().Value())
		}
	})
}

