package domain

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/jeffersonbrasilino/ddgo"
)

type Builder struct {
	buildErrors []string
	uuId        string
	username    string
	password    string
	person      *Person
}

type WithPersonProps struct {
	Person   *PersonProps
	Document *DocumentProps
	Contacts []*ContactProps
}

func NewBuilder() *Builder {
	return &Builder{
		buildErrors: make([]string, 0, 4),
	}
}

func (b *Builder) WithUuId(uuId string) *Builder {
	b.uuId = uuId
	return b
}

func (b *Builder) WithUsername(username string) *Builder {
	b.username = username
	return b
}

func (b *Builder) WithPassword(password string) *Builder {
	b.password = password
	return b
}

func (b *Builder) WithPerson(personProps *WithPersonProps) *Builder {

	if personProps.Person == nil {
		b.buildErrors = append(b.buildErrors, "person: person is required")
		return b
	}

	props := *personProps.Person

	errs := make([]string, 0, 2)
	if personProps.Document != nil {
		doc, err := NewDocument(personProps.Document)
		if err != nil {
			errs = append(errs, err.Error())
		}
		props.Document = doc
	}

	if personProps.Contacts != nil {
		contacts := make([]*Contact, 0, len(personProps.Contacts))
		for _, contact := range personProps.Contacts {
			contact, err := NewContact(contact)
			if err != nil {
				errs = append(errs, err.Error())
			}
			contacts = append(contacts, contact)
		}
		props.Contacts = contacts
	}

	if len(errs) > 0 {
		b.buildErrors = append(b.buildErrors, fmt.Sprintf("person: %s", strings.Join(errs, ", ")))
	}

	person, err := NewPerson(&props)
	if err != nil {
		b.buildErrors = append(b.buildErrors, err.Error())
		return b
	}

	b.person = person
	return b
}

func (b *Builder) Build() (*User, error) {
	if len(b.buildErrors) > 0 {
		validationResult, failed := json.Marshal(b.buildErrors)
		if failed != nil {
			return nil, ddgo.NewInternalError("Error when marshaling validation errors")
		}
		return nil, ddgo.NewInvalidDataError(string(validationResult))
	}

	return NewUser(&UserProps{
		UuId:     b.uuId,
		Username: b.username,
		Password: b.password,
		Person:   b.person,
	})
}
