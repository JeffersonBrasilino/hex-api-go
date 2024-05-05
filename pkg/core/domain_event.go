package core

import (
	uuidLib "github.com/google/uuid"

	"time"
)

type DomainEvent struct {
	uuid      string
	occuredOn time.Time
}

type domainEvent interface {
	GetPayload() any
	GetHeaders() any
}

func NewDomainEvent() DomainEvent {
	return DomainEvent{uuidLib.NewString(), time.Now()}
}

func (e *DomainEvent) GetOccuredOn() time.Time {
	return e.occuredOn
}

func (e *DomainEvent) GetUuid() string {
	return e.uuid
}
