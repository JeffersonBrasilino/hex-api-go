### Application layer guidelines

This document defines the application layer guidelines for the project. 
This layer acts as the orchestrator between the domain layer (business rules) and the infrastructure layer. It uses the CQRS (Command Query Responsibility Segregation) pattern to separate write operations (Commands) from read operations (Queries).

> **Note**: Currently, this guideline focuses on the **Command** side of the CQRS pattern. Query guidelines will be added once implemented.

#### Application layer structure

```
├── application/
│   ├── command/
│   │   └── [command_name]/
│   │       ├── command.go
│   │       └── handler.go
│   └── query/
│       └── [query_name]/
│           ├── query.go
│           └── handler.go
```

### Rules

- this layer is responsible for orchestrating domain actions. It acts as the "glue" between domain and infrastructure.
- must be independent of infrastructure details (e.g., HTTP requests, database transactions). It only knows about domain contracts.
- must follow the CQRS pattern, separating Commands (state-changing actions) from Queries (read-only actions).
- each action (Command/Query) must have its own dedicated directory.
- `command.go` must define the data structure of the command and implement a `Name()` method.
- `handler.go` must implement the `Handle` method to process the command, interacting with domain aggregates and domain contracts (repositories, gateways).
- the handler is responsible for creating/reconstituting the Domain Aggregate, performing operations on it, and persisting the state changes via domain contracts.
- handlers map external DTOs (from the command struct) into Domain objects.

### Name Conventions

- the action directory (e.g., `createuser`) must be lowercase, without underscores or spaces.
- within the action directory, the files must be named `command.go` and `handler.go`.
- the command struct must be named `Command` and the handler struct must be named `Handler`.
- the handler constructor must be named `NewComandHandler` (or `NewCommandHandler`).
- the command's `Name()` method must return the action name in camelCase (e.g., `createUser`).

### Implementation patterns

- Command -> `command_pattern.md`
- Command Handler -> `command_handler_pattern.md`
