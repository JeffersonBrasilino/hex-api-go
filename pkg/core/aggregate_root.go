package core

import "fmt"

type AggregateRoot struct {
	Entity
	domainEvents []domainEvent
}

func NewAggregateRoot(uuid string) AggregateRoot {
	return AggregateRoot{
		NewEntity(uuid),
		make([]domainEvent, 0),
	}
}

func (aggregate *AggregateRoot) GetUuid() string {
	return aggregate.Entity.GetUuid()
}

func (aggregate *AggregateRoot) GetDomainEvents() []domainEvent {
	return aggregate.domainEvents
}

func (aggregate *AggregateRoot) AddDomainEvent(event domainEvent) {
	fmt.Println("event added", event.GetHeaders(), event.GetPayload())
	//aggregate.domainEvents = append(aggregate.domainEvents, event)
}

func (aggregate *AggregateRoot) RemoveDomainEvent(event domainEvent) {

}