package domain

type AggregateRoot struct {
	Entity
	domainEvents []domainEvent
}

func NewAggregateRoot(uuid string) AggregateRoot {
	return AggregateRoot{
		Entity: NewEntity(uuid),
	}
}

func (aggregate *AggregateRoot) Uuid() string {
	return aggregate.Entity.Uuid()
}

func (aggregate *AggregateRoot) DomainEvents() []domainEvent {
	return aggregate.domainEvents
}

func (aggregate *AggregateRoot) AddDomainEvent(event domainEvent) {
	//fmt.Println("event added", event.Headers(), event.Payload())
	aggregate.domainEvents = append(aggregate.domainEvents, event)
}

func (aggregate *AggregateRoot) RemoveDomainEvent(event domainEvent) {

}