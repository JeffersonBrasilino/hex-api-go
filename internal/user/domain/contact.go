package domain

import (
	"encoding/json"

	"github.com/jeffersonbrasilino/ddgo"
)

type ContactProps struct {
	UuId        string `domainValidator:"required"`
	Description string `domainValidator:"required"`
	ContactType string `domainValidator:"required"`
}

type Contact struct {
	*ddgo.Entity
	uuId        string
	description string
	contactType string
}

func NewContact(props *ContactProps) (*Contact, error) {
	err := validateContact(props)
	if err != nil {
		return nil, err
	}
	return &Contact{
		description: props.Description,
		contactType: props.ContactType,
		Entity:      ddgo.NewEntity(props.UuId),
	}, nil
}

func validateContact(props *ContactProps) error {
	validator := ddgo.ValidatorInstance()
	validationErrors, faliedValidation := validator.Validate(props)
	if faliedValidation != nil {
		return ddgo.NewInternalError("Error when validating contact data")
	}

	if len(validationErrors) > 0 {
		validationResult, failed := json.Marshal(validationErrors)
		if failed != nil {
			return ddgo.NewInternalError("Error when marshaling validation errors")
		}
		return ddgo.NewInvalidDataError(string(validationResult))
	}

	return nil
}

func (c *Contact) Description() string {
	return c.description
}

func (c *Contact) ContactType() string {
	return c.contactType
}
