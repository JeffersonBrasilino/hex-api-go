// Package domain provides the core domain entities and value objects for the user bounded context.
// This file defines the Person entity, representing a natural person with personal data,
// contacts, and an associated document within the domain model.
package domain

import (
	"encoding/json"

	"github.com/jeffersonbrasilino/ddgo"
)

// PersonProps holds the input data required to create a Person entity.
// UuId, Name, and BirthDate are mandatory fields.
type PersonProps struct {
	UuId      string `domainValidator:"required"`
	Name      string `domainValidator:"required"`
	BirthDate string `domainValidator:"required"`
	Contacts  []*Contact
	Document  *Document
}

// Person is a domain entity that represents a natural person.
// It embeds ddgo.Entity for identity and lifecycle management.
type Person struct {
	*ddgo.Entity
	contacts  []*Contact
	document  *Document
	name      string
	birthDate string
}

// NewPerson creates and returns a new Person entity from the provided PersonProps.
// It validates the required fields using the domain validator before constructing the entity.
//
// Parameters:
//   - props: pointer to PersonProps containing the person's data.
//
// Returns a pointer to a valid Person and nil error on success, or nil and a domain error
// if validation fails.
func NewPerson(props *PersonProps) (*Person, error) {
	if err := validatePerson(props); err != nil {
		return nil, err
	}
	return &Person{
		name:      props.Name,
		birthDate: props.BirthDate,
		contacts:  props.Contacts,
		document:  props.Document,
		Entity:    ddgo.NewEntity(props.UuId),
	}, nil
}

// validatePerson runs structural validation against the given PersonProps using the domain validator.
//
// Parameters:
//   - props: pointer to PersonProps to validate.
//
// Returns nil on success, or a domain error describing the validation failure.
func validatePerson(props *PersonProps) error {
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

// Name returns the full name of the person.
func (p *Person) Name() string {
	return p.name
}

// Document returns the CPF document value object associated with the person, or nil if not set.
func (p *Person) Document() *Document {
	return p.document
}

// Contacts returns the list of contact entities associated with the person.
func (p *Person) Contacts() []*Contact {
	return p.contacts
}

// BirthDate returns the person's birth date string.
func (p *Person) BirthDate() string {
	return p.birthDate
}
