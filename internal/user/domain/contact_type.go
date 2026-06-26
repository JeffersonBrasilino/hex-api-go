// Package domain provides the core domain entities and value objects for the user bounded context.
// This file defines the ContactType entity, which categorizes contacts (e.g. email, phone)
// within the domain model.
package domain

import (
	"encoding/json"

	"github.com/jeffersonbrasilino/ddgo"
)

// ContactTypeProps holds the input data required to create a ContactType entity.
// UuId is the only mandatory field.
type ContactTypeProps struct {
	UuId        string `domainValidator:"required"`
	Description string
}

// ContactType is a domain entity that represents a category of contact information.
// It embeds ddgo.Entity for identity and lifecycle management.
type ContactType struct {
	*ddgo.Entity
	description string
}

// NewContactType creates and returns a new ContactType entity from the provided ContactTypeProps.
// It validates the required fields using the domain validator before constructing the entity.
//
// Parameters:
//   - props: pointer to ContactTypeProps containing the contact type's data.
//
// Returns a pointer to a valid ContactType and nil error on success, or nil and a domain error
// if validation fails.
func NewContactType(props *ContactTypeProps) (*ContactType, error) {
	if err := validateContactType(props); err != nil {
		return nil, err
	}
	return &ContactType{
		description: props.Description,
		Entity:      ddgo.NewEntity(props.UuId),
	}, nil
}

// validateContactType runs structural validation against the given ContactTypeProps using the domain validator.
//
// Parameters:
//   - props: pointer to ContactTypeProps to validate.
//
// Returns nil on success, or a domain error describing the validation failure.
func validateContactType(props *ContactTypeProps) error {
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

// Description returns the human-readable description of the contact type.
func (c *ContactType) Description() string {
	return c.description
}
