package domain

import (
	"encoding/json"

	"github.com/jeffersonbrasilino/ddgo"
)

type ContactTypeProps struct {
	UuId        string `domainValidator:"required"`
	Description string
}

type ContactType struct {
	*ddgo.Entity
	uuId        string
	description string
}

func NewContactType(props *ContactTypeProps) (*ContactType, error) {
	err := validateContactType(props)
	if err != nil {
		return nil, err
	}
	return &ContactType{
		description: props.Description,
		Entity:      ddgo.NewEntity(props.UuId),
	}, nil
}

func validateContactType(props *ContactTypeProps) error {
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

func (c *ContactType) Description() string {
	return c.description
}
