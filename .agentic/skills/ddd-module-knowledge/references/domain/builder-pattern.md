#### Aggregate Root builder pattern

The builder pattern is used to create new instances of an aggregate root, since aggregate roots can be composed of entities and therefore become complex to build.

The builder pattern should:

- Return an aggregate root
- Use the fluent interface pattern
- A single builder per aggregate root
- The builder should store creation errors and return them all at once, facilitating correction.

##### Gotchas
- If the aggregate root has properties of type entities or value objects, the function for creating this property can be complex. For these cases, use the `WithPerson` function of the builder `internal/user/domain/builder.go` as an example.

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

Implementation example: see -> `internal/user/domain/builder.go`
