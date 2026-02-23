package domain

/*
Abstração de raiz agregada.

Um agregado DDD é um cluster de objetos de domínio que pode ser tratado como uma única unidade.

[https://martinfowler.com/bliki/DDD_Aggregate.html]
*/
type AggregateRoot struct {
	*Entity
	domainEvents []domainEvent
}

// cria a intancia de raiz agregada
func NewAggregateRoot(uuid string) *AggregateRoot {
	return &AggregateRoot{
		NewEntity(uuid),
		[]domainEvent{},
	}
}

// retorna os eventos de dominios criados
func (a *AggregateRoot) DomainEvents() []domainEvent {
	return a.domainEvents
}

// adiciona um evento de dominio na lista de despacho de eventos
func (a *AggregateRoot) AddDomainEvent(event domainEvent) {
	a.domainEvents = append(a.domainEvents, event)
}

// remove um evento de dominio da lista de despacho
func (a *AggregateRoot) RemoveDomainEvent(event domainEvent) {
	for i, e := range a.domainEvents {
		if e == event {
			a.domainEvents = append(a.domainEvents[:i], a.domainEvents[i+1:]...)
			a.domainEvents[len(a.domainEvents)-1] = nil
			break
		}
	}
}

// esvazia a lista de despacho de eventos de dominio
func (a *AggregateRoot) ClearEvents() {
	a.domainEvents = []domainEvent{}
}
