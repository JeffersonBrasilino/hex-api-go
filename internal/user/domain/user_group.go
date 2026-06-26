// Package domain provides the core domain entities and value objects for the user bounded context.
// This file defines the UserGroup entity, which represents a named group of permissions
// that can be assigned to users within the domain model.
package domain

import (
	"encoding/json"

	"github.com/jeffersonbrasilino/ddgo"
)

// UserGroupProps holds the input data required to create a UserGroup entity.
// UuId is the only mandatory field.
type UserGroupProps struct {
	UuId        string `domainValidator:"required"`
	Name        string
	Permissions []any
	Main        bool
}

// UserGroup is a domain entity that represents a named permission group.
// It embeds ddgo.Entity for identity and lifecycle management.
type UserGroup struct {
	*ddgo.Entity
	name        string
	permissions []any
	main        bool
}

// NewUserGroup creates and returns a new UserGroup entity from the provided UserGroupProps.
// It validates the required fields using the domain validator before constructing the entity.
//
// Parameters:
//   - props: pointer to UserGroupProps containing the group's data.
//
// Returns a pointer to a valid UserGroup and nil error on success, or nil and a domain error
// if validation fails.
func NewUserGroup(props *UserGroupProps) (*UserGroup, error) {
	if err := validateUserGroup(props); err != nil {
		return nil, err
	}
	return &UserGroup{
		name:        props.Name,
		permissions: props.Permissions,
		main:        props.Main,
		Entity:      ddgo.NewEntity(props.UuId),
	}, nil
}

// validateUserGroup runs structural validation against the given UserGroupProps using the domain validator.
//
// Parameters:
//   - props: pointer to UserGroupProps to validate.
//
// Returns nil on success, or a domain error describing the validation failure.
func validateUserGroup(props *UserGroupProps) error {
	validator := ddgo.ValidatorInstance()
	validationErrors, err := validator.Validate(props)
	if err != nil {
		return ddgo.NewInternalError("Error when validating user group data")
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

// Name returns the display name of the user group.
func (ug *UserGroup) Name() string {
	return ug.name
}

// Permissions returns the list of permissions assigned to this user group.
func (ug *UserGroup) Permissions() []any {
	return ug.permissions
}

// Main reports whether this is the primary group for a user.
func (ug *UserGroup) Main() bool {
	return ug.main
}
