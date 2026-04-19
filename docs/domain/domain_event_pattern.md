#### Domain Event pattern

Domain events are objects that represent something that happened in the domain layer.

Example:

```go
package events

import "time"

type [event-name] struct {}

func New[event-name]() *[event-name] {
	return &[event-name]{}
}

func (e *[event-name]) Payload() any {
}

func (e *[event-name]) OcurredOn() time.Time {
	return time.Now()
}
```
Implementation example: see -> `../../internal/user/domain/events/user_created.go`
