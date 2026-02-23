package domain

import (
	"time"
)

type domainEvent interface {
	Payload() any
	Headers() any
}

/*
Abstração de evento de dominio
*/
type DomainEvent struct {
	uuid       string
	occurredOn time.Time
}


// cria instancia de evento de dominio
func NewDomainEvent(uuid string) DomainEvent {
	return DomainEvent{uuid, time.Now()}
}

// retorna a data que o evento aconteceu(criação de evento de dominio)
func (e *DomainEvent) OccurredOn() time.Time {
	return e.occurredOn
}

// retorna o UUID de evento de dominio
func (e *DomainEvent) Uuid() string {
	return e.uuid
}
