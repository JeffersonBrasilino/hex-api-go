// Package domain provides the core domain entities and value objects for the user bounded context.
// This file defines the User aggregate root, the central entity of this bounded context,
// which owns authentication credentials, personal data, and group assignments.
package domain

import (
	"encoding/json"

	"github.com/jeffersonbrasilino/ddgo"
)

// UserProps holds the input data required to create a User aggregate root.
// UuId, Username, Password, and Person are mandatory fields.
type UserProps struct {
	UuId       string    `domainValidator:"required"`
	Username   string    `domainValidator:"required,gte=1"`
	Password   *Password `domainValidator:"required"`
	Person     *Person   `domainValidator:"required"`
	UserGroups []*UserGroup
}

// User is the aggregate root for the user bounded context.
// It embeds ddgo.AggregateRoot for identity, lifecycle, and domain event management.
type User struct {
	*ddgo.AggregateRoot
	username   string
	password   *Password
	person     *Person
	userGroups []*UserGroup
}

// NewUser creates and returns a new User aggregate root from the provided UserProps.
// It validates all required fields using the domain validator before constructing the aggregate.
//
// Parameters:
//   - props: pointer to UserProps containing the user's data.
//
// Returns a pointer to a valid User and nil error on success, or nil and a domain error
// if validation fails.
func NewUser(props *UserProps) (*User, error) {
	if err := validate(props); err != nil {
		return nil, err
	}
	return &User{
		AggregateRoot: ddgo.NewAggregateRoot(props.UuId),
		username:      props.Username,
		password:      props.Password,
		person:        props.Person,
		userGroups:    props.UserGroups,
	}, nil
}

// validate runs structural validation against the given UserProps using the domain validator.
//
// Parameters:
//   - props: pointer to UserProps to validate.
//
// Returns nil on success, or a domain error describing the validation failure.
func validate(props *UserProps) error {
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

// Password returns the Password value object associated with the user.
func (u *User) Password() *Password {
	return u.password
}

// Username returns the login name of the user.
func (u *User) Username() string {
	return u.username
}

// SetPassword replaces the user's current password with the provided Password value object.
//
// Parameters:
//   - password: the new Password value object to assign.
func (u *User) SetPassword(password *Password) {
	u.password = password
}

// Person returns the Person entity associated with this user.
func (u *User) Person() *Person {
	return u.person
}

// UserGroups returns the list of UserGroup entities the user belongs to.
func (u *User) UserGroups() []*UserGroup {
	return u.userGroups
}
