package domain

import (
	"encoding/json"

	"github.com/jeffersonbrasilino/ddgo"
)

type DocumentProps struct {
	Value string `domainValidator:"required"`
}

type Document struct {
	value string
}

func NewDocument(props *DocumentProps) (*Document, error) {
	err := validateDocument(props)
	if err != nil {
		return nil, err
	}
	return &Document{
		value: props.Value,
	}, nil
}

func validateDocument(props *DocumentProps) error {
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

func (d *Document) Value() string {
	return d.value
}
