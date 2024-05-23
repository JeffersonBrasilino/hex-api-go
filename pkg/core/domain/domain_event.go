package domain

import (
	uuidLib "github.com/google/uuid"

	"time"
)

/*
Abstração de evento de dominio
*/
type DomainEvent struct {
	uuid       string
	occurredOn time.Time
}

type domainEvent interface {
	Payload() any
	Headers() any
}

// cria instancia de evento de dominio
func NewDomainEvent() DomainEvent {
	return DomainEvent{uuidLib.NewString(), time.Now()}
}

// retorna a data que o evento aconteceu(criação de evento de dominio)
func (e *DomainEvent) OccurredOn() time.Time {
	return e.occurredOn
}

// retorna o UUID de evento de dominio
func (e *DomainEvent) Uuid() string {
	return e.uuid
}
