package ddgo

// Entity is the base type for domain entities with a unique identifier.
//
// Intent: Provide a common abstraction where domain rules can live,
// independent of persistence or external integrations.
// Objective: Allow embedding Entity in concrete domain types so they
// share identity handling and stay free of infrastructure concerns.
//
// Example:
//
//	type item struct {
//	    *Entity
//	}
//	var itemInstance = &item{Entity: domain.NewEntity("")}
type Entity struct {
	uuid string
}

// NewEntity builds a new domain entity with the given UUID.
//
// Parameters:
//   - uuid: unique identifier for the entity (may be empty).
//
// Returns: a non-nil *Entity holding the given uuid.
func NewEntity(uuid string) *Entity {
	return &Entity{
		uuid: uuid,
	}
}

// Uuid returns the entity's unique identifier.
//
// Returns: the uuid string passed at construction (may be empty).
func (entity *Entity) Uuid() string {
	return entity.uuid
}
