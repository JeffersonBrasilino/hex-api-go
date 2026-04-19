### Infrastructure layer guidelines

This document defines the infrastructure layer guidelines for the project.
This layer is responsible for containing the implementation details of the module, such as database persistence, HTTP handlers, external service integrations (gateways), and messaging adapters.
All implementations in this layer must respect the contracts (interfaces) defined in the domain layer `contract` folder.

#### Infrastructure layer structure

```
├── infrastructure/
│   ├── database/
│   │   ├── gorm_model.go
│   │   ├── gorm_[module]_repository.go
│   │   └── mapper.go
│   └── http/
│       └── [action]_handler.go
```

### Rules

- this layer is responsible only for implementation details (adapters).
- must implement the contracts defined in the domain layer `contract` folder.
- must not contain business rules, only technical orchestration.
- must not expose domain internals to external frameworks (use mappers for data conversion).
- database models (persistence models) must be separate structs from domain entities.
- mapper functions must be package-private (`toDatabase`, `toDomain`), residing in the same package as the persistence layer.
- HTTP handlers must dispatch actions through the `gomes` bus (`CommandBus` or `QueryBus`), never calling domain or repository directly.
- repository implementations must use transaction management for write operations.
- repository errors must be wrapped with `ddgo` error types (`NewInternalError`, `NewInvalidDataError`).
- HTTP request structs must use `binding` tags for input validation.
- HTTP handlers must use `pkg/http` helpers for standardized responses (`http.Error`, `http.ErrorWithCode`, `http.Success`).
- OpenTelemetry tracing must be initialized per handler using `gomes/otel`.
- use the examples to base for implementation.

### Name Conventions

#### Database sub-layer

- Persistence model file name must be `gorm_model.go`.
- Persistence model struct name must be plural and PascalCase ex: `Users`, `PersonContacts`.
- Persistence model must define `TableName()` method with pattern `[project-name].[table_name]` ex: `hex-api-go.users`.
- Repository file name must follow `gorm_[module]_repository.go` ex: `gorm_user_repository.go`.
- Repository struct name must be prefixed with ORM name and PascalCase ex: `GormUserRepository`.
- Repository constructor must follow `New[OrmName][Module]Repository` ex: `NewGormUserRepository`.
- Mapper file name must be `mapper.go`.
- Mapper functions must be package-private and named `toDatabase` and `toDomain`.

#### HTTP sub-layer

- Handler file name must follow `[action_name]_handler.go` using snake_case ex: `create_user_handler.go`.
- Handler function name must be PascalCase ex: `CreateUserHandler`.
- Request struct name must follow `[ActionName]Request` ex: `CreateUserRequest`.
- Trace variable must follow `[actionName]Trace` ex: `createUserTrace`.

### Implementation patterns

- Repository Implementation -> `repository_pattern.md`
- Persistence Models (GORM) -> `persistence_model_pattern.md`
- Domain-Database Mapper -> `mapper_pattern.md`
- HTTP Handler -> `http_handler_pattern.md`
