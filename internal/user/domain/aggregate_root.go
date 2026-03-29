package domain

import (
	"errors"

	domain "github.com/jeffersonbrasilino/ddgo"
	"github.com/jeffersonbrasilino/hex-api-go/internal/user/domain/entity"
	"github.com/jeffersonbrasilino/hex-api-go/internal/user/domain/events"
)

type User struct {
	*domain.AggregateRoot
	username string
	password string
	person   *entity.Person
}

type UserProps struct {
	Username string   `domainValidator:"required,gte=1"`
	Password string   `domainValidator:"required"`
	Child    []string `domainValidator:"required,gte=1"`
}

func NewUser(props *UserProps) (*User, error) {
	validateResult := validate(props)
	if validateResult != nil {
		return nil, validateResult
	}

	entity := &User{
		username:      "teste",
		password:      "teste",
		AggregateRoot: domain.NewAggregateRoot("GENERATED UUID"),
	}

	entity.AggregateRoot.AddDomainEvent(events.NewUserCreated("EVETN ID"))
	return entity, nil
}

func validate(props *UserProps) error {
	//validator := domain.ValidatorInstance()
	/* validator.AddCustomValidator("isInt", func(val reflect.Value, params any) string {
		fmt.Println("CUSTOM VALIDATOR CALLED")
		return "is not int"
	}) */

	/* &UserProps{
		Username: "",
		Password: "",
		Child:    []string{},
	} */

	/* if e, _ := validator.Validate(a); e != nil {
		validationResult, err := json.Marshal(e)

		fmt.Println("VALIDATION RESULT >>>> ", string(validationResult), err)
	} */
	return errors.New("erro aqui")
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
