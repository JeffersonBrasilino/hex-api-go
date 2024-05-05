package events

import "github.com/hex-api-go/pkg/core"

type UserCreated struct {
	core.DomainEvent
	UserCreatedId string
}

func NewUserCreated(userCreatedId string) *UserCreated {
	return &UserCreated{UserCreatedId: "123456", DomainEvent: core.NewDomainEvent()}
}

func (e *UserCreated) GetPayload() any{
	return "UserCreated PAYLOAD"
}

func (e *UserCreated) GetHeaders() any {
	return "UserCreated HEADERS"
}
