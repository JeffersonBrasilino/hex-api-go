### Domain layer guidelines

This document defines the domain layer guidelines for the project. This layer usages ddgo plugin.
This layer is responsable only to contain business rules, For implementation details, this layer will contain the respective contracts(repository, gateways, etc...) defined in the `contract` folder

#### Domain layer structure

```
├── domain/
│   ├── contract/
│   │   ├── repository.go
│   │   └── service.go
│   ├── event/
│   │   └── event.go
│   ├── entity.go
│   └── builder.go
```

### Rules

- this layer is responsable only to contain business rules.
- must be independent of the other layers
- must not have any dependencies on other layers
- Aggregate root must be the only entry point to the domain
- Aggregate root mus be builder pattern to create new instances
- Domain Contracts must be implemented in the `infrastructure` folder
- use the examples to base for implementation
- domain events must be implemented in the `event` folder
- domain contracts must be implemented in the `contract` folder
- only aggregate root makes domain events

### Name Conventions

- Entity and value objects file name must be singular and snake_case ex: `user.go`, `person_contact.go`.
- Entity and value objects struct name must be singular and PascalCase ex: `User`, `PersonContact`.
- builder file name mus be `builder.go`.
- entity aggregate root must be name equal to module name.
- domain events file name must be singular and snake_case ex: `user_created.go`.
- domain events struct name must be singular and PascalCase ex: `UserCreated`.
- domain contracts file name must be singular and snake_case ex: `user_repository.go`.
- domain contracts struct name must be singular and PascalCase ex: `UserRepository`.

This layer components use ddgo plugin, see https://github.com/jeffersonbrasilino/ddgo for more information.

### Implementation patterns

- Entity or Aggregate Root -> `entity_pattern.md`
- Aggregate Root Builder -> `builder_pattern.md`
- Value Objects -> `value_object_pattern.md`
- Domain Events -> `domain_event_pattern.md`
- Domain Contracts -> `contract_pattern.md`
