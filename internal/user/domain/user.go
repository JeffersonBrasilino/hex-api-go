package domain

import (
	"github.com/hex-api-go/internal/user/domain/entity"
	"github.com/hex-api-go/internal/user/domain/events"
	"github.com/hex-api-go/pkg/core"
)

type User struct {
	core.AggregateRoot
	username string
	password string
	person   *entity.Person
}

func NewUser(username string, password string) *User {
	entity := &User{username: username, password: password, AggregateRoot: core.NewAggregateRoot("GENERATED UUID")}
	entity.AggregateRoot.AddDomainEvent(events.NewUserCreated("EVETN ID"))
	return entity
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
