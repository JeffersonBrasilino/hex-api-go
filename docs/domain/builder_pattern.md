#### Aggregate Root builder pattern

Aggregate Root must be a builder pattern to create new instances.
ONLY aggregate root entity must be a builder pattern to create new instances.
The builder creates the child entities for the aggregated root.

Boilerplate Example:

```go
package domain

import (
	"encoding/json"
	"fmt"
	"strings"
	"github.com/jeffersonbrasilino/ddgo"
)

type Builder struct {
	prop1 string
	prop2 string
	childEntity *child-entity //only if aggregate depends on other entities
	buildErrors []string
}

//when aggregate depends on other entities, use this struct to receive child entity props. 
//only if aggregate depends on other entities
type With[child-entity]Props struct {
}

func NewBuilder() *Builder {
	return &Builder{
		buildErrors: make([]string, 0, 4),
	}
}

func (b *Builder) WithProp1(prop1 string) *Builder {
	b.prop1 = prop1
	return b
}

func (b *Builder) WithProp2(prop2 string) *Builder {
	b.prop2 = prop2
	return b
}

//this method use With[child-entity]Props to create child entity
//only if aggregate depends on other entities
func (b *Builder) With[child-entity](props *With[child-entity]Props) *Builder {
	// ...build child entity
}

func (b *Builder) Build() (*User, error) {
	if len(b.buildErrors) > 0 {
		validationResult, failed := json.Marshal(b.buildErrors)
		if failed != nil {
			return nil, ddgo.NewInternalError("Error when marshaling validation errors")
		}
		return nil, ddgo.NewInvalidDataError(string(validationResult))
	}

	return New[entity-name](&[entity-name]Props{
		Prop1: b.prop1,
		Prop2: b.prop2,
		ChildEntity: b.childEntity,
	})
}
```
Implementation example: see -> `../../internal/user/domain/builder.go`
