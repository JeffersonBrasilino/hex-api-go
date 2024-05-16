package domain

import (
	"time"

	"github.com/hex-api-go/internal/user/domain/entity"
	"github.com/hex-api-go/internal/user/domain/events"
	"github.com/hex-api-go/pkg/core/domain"
)

type User struct {
	domain.AggregateRoot
	username string
	password string
	person   *entity.Person
}

func NewUser(username string, password string) *User {
	entity := &User{username: username, password: password, AggregateRoot: domain.NewAggregateRoot("GENERATED UUID")}
	entity.AggregateRoot.AddDomainEvent(events.NewUserCreated("EVETN ID"))
	time.Sleep(time.Second * 5)
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
