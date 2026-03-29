package events

import "time"

type UserCreated struct {
	UserCreatedId string
}

func NewUserCreated(userCreatedId string) *UserCreated {
	return &UserCreated{UserCreatedId: "123456"}
}

func (e *UserCreated) Payload() any {
	return "UserCreated PAYLOAD"
}

func (e *UserCreated) OcurredOn() time.Time {
	return time.Now()
}

func (e *UserCreated) Uuid() string {
	return "hueheuhueh"
}
