#### Entity or Aggregate Root pattern

Entity or Aggregate Root is a object that has a identity and can be mutated over time.
It is the main building block of the domain layer.
The entity may contain value objects or other entities.

Boilerplate Example:

```go
package domain

import (
	"encoding/json"
	"github.com/jeffersonbrasilino/ddgo"
)

type [entity-name]Props struct {
	UuId     string         `domainValidator:"required"`
}

type [entity-name] struct {
	*ddgo.AggregateRoot // for aggregate root entity, for entity use *ddgo.Entity
}

func New[entity-name](props *[entity-name]Props) (*[entity-name], error) {
	validateResult := validate(props)
	if validateResult != nil {
		return nil, validateResult
	}

	return &[entity-name]{
		AggregateRoot: ddgo.NewAggregateRoot(props.UuId), // for entity use *ddgo.NewEntity
	}, nil
}

func validate(props *[entity-name]Props) error {
	validator := ddgo.ValidatorInstance()
	validationErrors, faliedValidation := validator.Validate(props)
	if faliedValidation != nil {
		return ddgo.NewInternalError("Error when validating contact data")
	}

	if len(validationErrors) > 0 {
		validationResult, failed := json.Marshal(validationErrors)
		if failed != nil {
			return ddgo.NewInternalError("Error when marshaling validation errors")
		}
		return ddgo.NewInvalidDataError(string(validationResult))
	}

	return nil
}
```

entity/aggregate root implementation example:

- for entity see -> `internal/user/domain/person.go`
- for aggregate root see -> `internal/user/domain/user.go`
