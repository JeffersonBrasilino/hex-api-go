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

//cria a intancia de raiz agregada
func NewAggregateRoot(uuid string) *AggregateRoot {
	return &AggregateRoot{
		Entity: NewEntity(uuid),
	}
}

//retorna o UUID da Entidade
func (aggregate *AggregateRoot) Uuid() string {
	return aggregate.Entity.Uuid()
}
//retorna os eventos de dominios criados
func (aggregate *AggregateRoot) DomainEvents() []domainEvent {
	return aggregate.domainEvents
}

//adiciona um evento de dominio na lista de despacho de eventos
func (aggregate *AggregateRoot) AddDomainEvent(event domainEvent) {
	aggregate.domainEvents = append(aggregate.domainEvents, event)
}

// remove um evento de dominio da lista de despacho
func (aggregate *AggregateRoot) RemoveDomainEvent(event domainEvent) {

}

//esvazia a lista de despacho de eventos de dominio
func (aggregate *AggregateRoot) ClearEvents(){

}
