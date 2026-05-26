package domain

import (
	"encoding/json"

	"github.com/jeffersonbrasilino/ddgo"
)

type UserGroupProps struct {
	UuId        string `domainValidator:"required"`
	Name        string
	Permissions []any
	Main        bool
}

type UserGroup struct {
	*ddgo.Entity
	name        string
	permissions []any
	main        bool
}

func NewUserGroup(props *UserGroupProps) (*UserGroup, error) {
	err := validateUserGroup(props)
	if err != nil {
		return nil, err
	}
	return &UserGroup{
		name:        props.Name,
		permissions: props.Permissions,
		main:        props.Main,
		Entity:      ddgo.NewEntity(props.UuId),
	}, nil
}

func validateUserGroup(props *UserGroupProps) error {
	validator := ddgo.ValidatorInstance()
	validationErrors, faliedValidation := validator.Validate(props)
	if faliedValidation != nil {
		return ddgo.NewInternalError("Error when validating user group data")
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

func (ug *UserGroup) Name() string {
	return ug.name
}

func (ug *UserGroup) Permissions() []any {
	return ug.permissions
}

func (ug *UserGroup) Main() bool {
	return ug.main
}
