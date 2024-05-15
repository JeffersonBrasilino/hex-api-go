package domain

import (
	uuidLib "github.com/google/uuid"

	"time"
)

type DomainEvent struct {
	uuid      string
	occuredOn time.Time
}

type domainEvent interface {
	Payload() any
	Headers() any
}

func NewDomainEvent() DomainEvent {
	return DomainEvent{uuidLib.NewString(), time.Now()}
}

func (e *DomainEvent) OccuredOn() time.Time {
	return e.occuredOn
}

func (e *DomainEvent) Uuid() string {
	return e.uuid
}
