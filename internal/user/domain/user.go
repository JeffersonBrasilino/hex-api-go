package domain

import (
	"encoding/json"
	"errors"

	"github.com/hex-api-go/internal/user/domain/entity"
	"github.com/hex-api-go/internal/user/domain/events"
	"github.com/hex-api-go/pkg/core/domain"
	domainValidator "github.com/hex-api-go/pkg/core/domain/validator"
)

type User struct {
	*domain.AggregateRoot
	username string
	password string
	person   *entity.Person
}

type UserProps struct {
	Username string `domainValidator:"required"`
	Password string `domainValidator:"required"`
}

func NewUser(username string, password string) (*User, error) {
	props := &UserProps{username, password}
	validateResult := validate(props)
	if validateResult != nil {
		return nil, validateResult
	}

	entity := &User{username: username, password: password, AggregateRoot: domain.NewAggregateRoot("GENERATED UUID")}
	entity.AggregateRoot.AddDomainEvent(events.NewUserCreated("EVETN ID"))
	return entity, nil
}

func validate(props *UserProps) error {
	validator := domainValidator.NewDomainValidator()
	/* validator.AddCustomValidator("isInt", func(val reflect.Value, params any) string {
		fmt.Println("CUSTOM VALIDATOR CALLED")
		return "is not int"
	}) */

	if !validator.Validate(props) {
		errString, _ := json.Marshal(validator.GetErrors())
		return errors.New(string(errString))
	}
	return nil
}

func (u *User) GetPassword() string {
	return u.password
}

func (u *User) GetUsername() string {
	return u.username
}

func (u *User) SetPassword(password string) {
	u.password = password
}

func (u *User) GetPerson() *entity.Person {
	return u.person
}
