#### Domain-Database Mapper pattern

Mapper functions are responsible for converting between domain aggregates/entities and persistence models.
They act as an anti-corruption layer, preventing database details from leaking into the domain.

The mapper must:

- be defined in a `mapper.go` file inside the `database` package.
- contain package-private functions only (unexported, lowercase).
- implement two directional functions: `toDatabase` (domain -> persistence) and `toDomain` (persistence -> domain).
- access domain entity data through getter methods only (respecting encapsulation).
- build persistence models by mapping each domain field to the corresponding model field.

Boilerplate Example:

```go
package database

import "github.com/jeffersonbrasilino/hex-api-go/internal/[module-name]/domain"

func toDomain(model *[PersistenceModel]) *domain.[AggregateRoot] {
	return &domain.[AggregateRoot]{}
}

func toDatabase(aggregate *domain.[AggregateRoot]) *[PersistenceModel] {
	return &[PersistenceModel]{
		Field1: aggregate.Field1(),
		Field2: aggregate.Field2(),
		Child: [ChildModel]{
			Name:  aggregate.Child().Name(),
			Value: aggregate.Child().ValueObject().Value(),
		},
	}
}
```

Implementation example: see -> `internal/user/infrastructure/database/mapper.go`
