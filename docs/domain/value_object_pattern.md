#### Value Objects pattern

Value Objects are objects that do not have a identity and are immutable over time.
It is the main building block of the domain layer.
The entity may contain value objects or other entities.

Example:

```go
package domain

import (
	"encoding/json"

	"github.com/jeffersonbrasilino/ddgo"
)

type [value-object-name]Props struct {
	Value string `domainValidator:"required"`
}

type [value-object-name] struct {
	value string
}

func New[value-object-name](props *[value-object-name]Props) (*[value-object-name], error) {
	err := validate[value-object-name](props)
	if err != nil {
		return nil, err
	}
	return &[value-object-name]{
		value: props.Value,
	}, nil
}

func validate[value-object-name](props *[value-object-name]Props) error {
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

func (d *[value-object-name]) Value() string {
	return d.value
}
```
Implementation example: see -> `../../internal/user/domain/document.go`
