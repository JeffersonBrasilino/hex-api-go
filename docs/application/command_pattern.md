#### Command pattern

The Command represents an intent to change the system's state. It is a simple Data Transfer Object (DTO) that holds all the necessary information to execute a specific business use case.
Commands should be immutable and not contain any business logic.

The command must:
- be defined in a `command.go` file inside its specific action directory.
- be a struct named `Command`.
- contain fields mapping exactly to what the use case requires to be executed (often decorated with `json` tags if originating from payload deserialization, though application layer should ideally be input-agnostic).
- implement a `Name()` method that returns a string identifier for the command (typically camelCase).

Boilerplate Example:

```go
package [actionname]

type Command struct {
	Field1     string `json:"field1"`
	Field2     string `json:"field2"`
	ChildValue string `json:"childValue"`
}

func (c *Command) Name() string {
	return "[actionName]"
}
```
Implementation example: see -> `../../internal/user/application/command/createuser/command.go`
