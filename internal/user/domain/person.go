package domain

import (
	"encoding/json"

	"github.com/jeffersonbrasilino/ddgo"
)

type PersonProps struct {
	UuId      string `domainValidator:"required"`
	Name      string `domainValidator:"required"`
	BirthDate string `domainValidator:"required"`
	Contacts  []*Contact
	Document  *Document
}

type Person struct {
	*ddgo.Entity
	contacts  []*Contact
	document  *Document
	name      string
	birthDate string
}

func NewPerson(props *PersonProps) (*Person, error) {
	err := validatePerson(props)
	if err != nil {
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

func validatePerson(props *PersonProps) error {
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

func (p *Person) Name() string {
	return p.name
}

func (p *Person) Document() *Document {
	return p.document
}

func (p *Person) Contacts() []*Contact {
	return p.contacts
}

func (p *Person) BirthDate() string {
	return p.birthDate
}
