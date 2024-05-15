package events

import "github.com/hex-api-go/pkg/core/domain"

type UserCreated struct {
	domain.DomainEvent
	UserCreatedId string
}

func NewUserCreated(userCreatedId string) *UserCreated {
	return &UserCreated{UserCreatedId: "123456", DomainEvent: domain.NewDomainEvent()}
}

func (e *UserCreated) Payload() any{
	return "UserCreated PAYLOAD"
}

func (e *UserCreated) Headers() any {
	return "UserCreated HEADERS"
}
