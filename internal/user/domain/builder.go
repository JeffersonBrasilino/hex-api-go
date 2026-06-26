// Package domain provides the core domain entities and value objects for the user bounded context.
// This file defines the Builder, which implements the Builder pattern for constructing a
// fully assembled User aggregate root from raw input data.
package domain

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/jeffersonbrasilino/ddgo"
)

// Builder accumulates raw input data and validation errors to construct a User aggregate.
// Use NewBuilder to obtain an instance, chain the With* methods, and call Build to finalize.
type Builder struct {
	buildErrors []string
	uuId        string
	username    string
	password    *Password
	person      *Person
	userGroups  []*UserGroup
}

// WithContactProps holds the raw data for a single contact entry used within WithPerson.
type WithContactProps struct {
	UuId        string
	Description string
	Main        bool
	ContactType string
}

// WithPersonProps holds the raw data for building a Person entity inside the Builder.
type WithPersonProps struct {
	UuId      string
	Name      string
	BirthDate string
	Document  string
	Contacts  []*WithContactProps
}

// NewBuilder creates and returns an empty Builder ready to receive input via the With* methods.
func NewBuilder() *Builder {
	return &Builder{
		buildErrors: make([]string, 0, 4),
	}
}

// WithUuId sets the unique identifier for the User being built.
//
// Parameters:
//   - uuId: the UUID string to assign to the user.
func (b *Builder) WithUuId(uuId string) *Builder {
	b.uuId = uuId
	return b
}

// WithUsername sets the username for the User being built.
//
// Parameters:
//   - username: the login name string to assign to the user.
func (b *Builder) WithUsername(username string) *Builder {
	b.username = username
	return b
}

// WithPassword creates a Password value object from the raw string and assigns it to the builder.
// If password validation fails, the error is collected and Build will return it.
//
// Parameters:
//   - password: the raw password string to validate and assign.
func (b *Builder) WithPassword(password string) *Builder {
	pwd, err := NewPassword(&PasswordProps{Value: password})
	if err != nil {
		b.buildErrors = append(b.buildErrors, err.Error())
		return b
	}
	b.password = pwd
	return b
}

// WithPerson builds a Person entity (including its Document and Contacts) from raw props
// and assigns it to the builder. Any validation errors encountered are collected and
// Build will return them.
//
// Parameters:
//   - personProps: pointer to WithPersonProps containing the person's raw data.
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

			c, err := NewContact(&ContactProps{
				UuId:        contact.UuId,
				Description: contact.Description,
				Main:        contact.Main,
				ContactType: contactType,
			})
			if err != nil {
				errs = append(errs, err.Error())
				continue
			}
			personEntityProps.Contacts = append(personEntityProps.Contacts, c)
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

// WithUserGroups builds UserGroup entities from the provided props slice and assigns them
// to the builder. Any validation errors encountered are collected and Build will return them.
//
// Parameters:
//   - groups: slice of pointers to UserGroupProps; a nil slice is a no-op.
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

// Build finalizes the construction of the User aggregate root.
// If any errors were collected during the With* calls, it returns a domain validation error
// containing all accumulated messages. Otherwise it delegates to NewUser.
//
// Returns a pointer to a valid User and nil error on success, or nil and a domain error
// if any accumulated or structural validation failed.
func (b *Builder) Build() (*User, error) {
	if len(b.buildErrors) > 0 {
		validationResult, err := json.Marshal(b.buildErrors)
		if err != nil {
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
