#### Command Handler pattern

The Command Handler is responsible for orchestrating the execution of a Command. It coordinates the Domain objects (Entities, Aggregates) and infrastructural components (via Domain Contracts) to fulfill a use case.

The command handler must:

- be defined in a `handler.go` file inside its specific action directory.
- be a struct named `Handler`.
- receive its dependencies (domain contracts like repositories, gateways) via constructor injection.
- implement a `Handle(ctx context.Context, data *Command) (any, error)` method.
- translate the data from the `Command` DTO into an Aggregate Root using the appropriate domain builder.
- delegate the actual business logic to the Aggregate Root if applicable, or manage the flow.
- use injected domain contracts (e.g., repository) to persist changes.
- not contain infrastructural details like HTTP contexts or direct database queries.

Boilerplate Example:

```go
package [actionname]

import (
	"context"

	"github.com/jeffersonbrasilino/gomes/otel"
	"github.com/jeffersonbrasilino/hex-api-go/internal/[module-name]/domain"
	"github.com/jeffersonbrasilino/hex-api-go/internal/[module-name]/domain/contract"
)

type Handler struct {
	repository    contract.[AggregateRoot]Repository
	tracer        otel.OtelTrace
	messageHeader map[string]string
}

func NewCommandHandler(repository contract.[AggregateRoot]Repository) *Handler {
	return &Handler{
		repository: repository,
	}
}

func (c *Handler) Handle(ctx context.Context, data *Command) (any, error) {
	// 1. Map Command DTO to Domain Aggregate
	aggregate, errAg := c.makeAggregate(data)
	if errAg != nil {
		return nil, errAg
	}

	// 2. Perform actions on aggregate or use contracts to persist
	err := c.repository.Create(ctx, aggregate)
	if err != nil {
		return nil, err
	}

	// 3. Return result or dispatch events if applicable
	return "success", nil
}

// helper method to encapsulate domain builder complexity
func (c *Handler) makeAggregate(data *Command) (*domain.[AggregateRoot], error) {
	return domain.NewBuilder().
		WithField1(data.Field1).
		WithField2(data.Field2).
		// map other fields
		Build()
}
```

Implementation example: see -> `internal/user/application/command/createuser/handler.go`
