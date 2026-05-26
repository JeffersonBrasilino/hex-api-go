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
	password    *Password
	person      *Person
	userGroups  []*UserGroup
}

type WithContactProps struct {
	UuId        string
	Description string
	Main        bool
	ContactType string
}

type WithPersonProps struct {
	UuId      string
	Name      string
	BirthDate string
	Document  string
	Contacts  []*WithContactProps
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
	pwd, err := NewPassword(&PasswordProps{Value: password})
	if err != nil {
		b.buildErrors = append(b.buildErrors, err.Error())
		return b
	}
	b.password = pwd
	return b
}

func (b *Builder) WithPerson(personProps *WithPersonProps) *Builder {

	personEntityProps := &PersonProps{
		UuId:      personProps.UuId,
		Name:      personProps.Name,
		BirthDate: personProps.BirthDate,
		Contacts:  make([]*Contact, 0, len(personProps.Contacts)),
	}
	errs := make([]string, 0, 2)
	if personProps.Document != "" {
		doc, err := NewDocument(&DocumentProps{Value: personProps.Document})
		if err != nil {
			errs = append(errs, err.Error())
			return b
		}
		personEntityProps.Document = doc
	}

	if personProps.Contacts != nil {
		for _, contact := range personProps.Contacts {
			if contact.ContactType == "" {
				errs = append(errs, "contact type is required")
				continue
			}

			contactType, err := NewContactType(&ContactTypeProps{UuId: contact.ContactType})
			if err != nil {
				errs = append(errs, err.Error())
				continue
			}
			contactEntityProps := &ContactProps{
				UuId:        contact.UuId,
				Description: contact.Description,
				Main:        contact.Main,
				ContactType: contactType,
			}
			contact, err := NewContact(contactEntityProps)

			if err != nil {
				errs = append(errs, err.Error())
				continue
			}
			personEntityProps.Contacts = append(personEntityProps.Contacts, contact)
		}
	}

	if len(errs) > 0 {
		b.buildErrors = append(b.buildErrors, fmt.Sprintf("person: %s", strings.Join(errs, ", ")))
	}

	person, err := NewPerson(personEntityProps)
	if err != nil {
		b.buildErrors = append(b.buildErrors, err.Error())
		return b
	}

	b.person = person
	return b
}

func (b *Builder) WithUserGroups(groups []*UserGroupProps) *Builder {
	if groups == nil {
		return b
	}

	userGroups := make([]*UserGroup, 0, len(groups))
	for _, groupProps := range groups {
		group, err := NewUserGroup(groupProps)
		if err != nil {
			b.buildErrors = append(b.buildErrors, err.Error())
			continue
		}
		userGroups = append(userGroups, group)
	}

	b.userGroups = userGroups
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
		UuId:       b.uuId,
		Username:   b.username,
		Password:   b.password,
		Person:     b.person,
		UserGroups: b.userGroups,
	})
}
