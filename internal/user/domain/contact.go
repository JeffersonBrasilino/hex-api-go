// Package domain provides the core domain entities and value objects for the user bounded context.
// This file defines the Contact entity, which represents a contact channel (e.g. phone number,
// email address) linked to a person within the domain model.
package domain

import (
	"encoding/json"

	"github.com/jeffersonbrasilino/ddgo"
)

// ContactProps holds the input data required to create a Contact entity.
// UuId, Description, and ContactType are mandatory fields.
type ContactProps struct {
	UuId        string       `domainValidator:"required"`
	Description string       `domainValidator:"required"`
	ContactType *ContactType `domainValidator:"required"`
	Main        bool
}

// Contact is a domain entity that represents a contact channel for a person.
// It embeds ddgo.Entity for identity and lifecycle management.
type Contact struct {
	*ddgo.Entity
	description string
	contactType *ContactType
	main        bool
}

// NewContact creates and returns a new Contact entity from the provided ContactProps.
// It validates the required fields using the domain validator before constructing the entity.
//
// Parameters:
//   - props: pointer to ContactProps containing the contact's data.
//
// Returns a pointer to a valid Contact and nil error on success, or nil and a domain error
// if validation fails.
func NewContact(props *ContactProps) (*Contact, error) {
	if err := validateContact(props); err != nil {
		return nil, err
	}
	return &Contact{
		description: props.Description,
		contactType: props.ContactType,
		main:        props.Main,
		Entity:      ddgo.NewEntity(props.UuId),
	}, nil
}

// validateContact runs structural validation against the given ContactProps using the domain validator.
//
// Parameters:
//   - props: pointer to ContactProps to validate.
//
// Returns nil on success, or a domain error describing the validation failure.
func validateContact(props *ContactProps) error {
	validator := ddgo.ValidatorInstance()
	validationErrors, err := validator.Validate(props)
	if err != nil {
		return ddgo.NewInternalError("Error when validating contact data")
	}

	if len(validationErrors) > 0 {
		validationResult, err := json.Marshal(validationErrors)
		if err != nil {
			return ddgo.NewInternalError("Error when marshaling validation errors")
		}
		return ddgo.NewInvalidDataError(string(validationResult))
	}

	return nil
}

// Description returns the contact value string (e.g. phone number or email address).
func (c *Contact) Description() string {
	return c.description
}

// ContactType returns the ContactType entity that categorizes this contact.
func (c *Contact) ContactType() *ContactType {
	return c.contactType
}

// IsMain reports whether this contact is the primary contact for its owner.
func (c *Contact) IsMain() bool {
	return c.main
}
