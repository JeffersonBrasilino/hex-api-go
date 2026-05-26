package domain

import (
	"encoding/json"

	"github.com/jeffersonbrasilino/ddgo"
	domain "github.com/jeffersonbrasilino/ddgo"
)

type UserProps struct {
	UuId       string    `domainValidator:"required"`
	Username   string    `domainValidator:"required,gte=1"`
	Password   *Password `domainValidator:"required"`
	Person     *Person   `domainValidator:"required"`
	UserGroups []*UserGroup
}

type User struct {
	*domain.AggregateRoot
	username   string
	password   *Password
	person     *Person
	userGroups []*UserGroup
}

func NewUser(props *UserProps) (*User, error) {
	validateResult := validate(props)
	if validateResult != nil {
		return nil, validateResult
	}

	entity := &User{
		AggregateRoot: domain.NewAggregateRoot(props.UuId),
		username:      props.Username,
		password:      props.Password,
		person:        props.Person,
		userGroups:    props.UserGroups,
	}

	return entity, nil
}

func validate(props *UserProps) error {
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

func (u *User) Password() *Password {
	return u.password
}

func (u *User) Username() string {
	return u.username
}

func (u *User) SetPassword(password *Password) {
	u.password = password
}

func (u *User) Person() *Person {
	return u.person
}

func (u *User) UserGroups() []*UserGroup {
	return u.userGroups
}
