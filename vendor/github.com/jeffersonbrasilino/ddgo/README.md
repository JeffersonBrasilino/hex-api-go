## ddgo – Domain-Driven Go Core

`ddgo` is a small, opinionated toolkit for building **Domain-Driven Design (DDD)** applications
in Go. It provides the core building blocks for rich domain models:

- **Entities**: domain objects with identity
- **Aggregate roots** and **domain events**
- **Tag-based validation** for DTOs and domain inputs
- **Standardized domain error types**

The library is intentionally infrastructure-agnostic. It can be used on its own in a clean
architecture/hexagonal architecture, or integrated with a message system such as the
`gomes` plugin (for CQRS, Kafka, RabbitMQ, etc.).

## 1. Overview

**Intent**: Provide a minimal, reusable domain core for Go projects that follow DDD, CQRS and
event-driven architecture.

**Objective**: Centralize the most common domain patterns so each project does not have to
re-implement:

- Entity identity and aggregate roots
- Domain events and event buffering
- Domain-level validation rules
- Well-known domain error types

### 1.1. Key Features

- **Pure domain layer**: no dependency on HTTP, databases, or message brokers.
- **Rich domain model**: entities and aggregate roots hold behavior, not just data.
- **Domain events**: aggregates can record events for later publishing.
- **Tag-based validation**: struct tags drive validation rules, keeping DTOs explicit.
- **Standard errors**: consistent error types like `NotFoundError`, `ValidationError`, etc.

### 1.2. Folder Structure

Inside `pkg/core/ddgo`:

- `entity.go`: base `Entity` type.
- `aggregate_root.go`: `AggregateRoot` and `DomainEvent` interface.
- `domain_validator.go`: validator entry point and reflection-based schema builder.
- `validators.go`: concrete validators (`required`, `gte`, `lte`, `len`).
- `errors.go`: domain error hierarchy.
- `*_test.go`: unit tests showing usage of each component.

---

## 2. Bootstrap

`ddgo` is a pure library – there is **no runtime to start or stop**. Bootstrapping mostly
means:

- Adding the module to your `go.mod`.
- Designing your domain model to embed `Entity` and `AggregateRoot`.
- Using the `ValidatorInstance` in your application/service layer.

### 2.1. Installing and Importing

In your project:

```bash
go get github.com/jeffersonbrasilino/ddgo
```

In your domain package:

```go
package user

import "github.com/jeffersonbrasilino/ddgo"
```

### 2.2. Registering Domain Components

Entities and aggregate roots are **just Go types** that embed `Entity` or `AggregateRoot`:

```go
package user

import "github.com/jeffersonbrasilino/ddgo"

type User struct {
    *domain.Entity
    Name string
}

func NewUser(id, name string) *User {
    return &User{
        Entity: domain.NewEntity(id),
        Name:   name,
    }
}
```

Validation is configured using struct tags and the shared `ValidatorInstance`:

```go
type CreateUserDTO struct {
    Name  string `domainValidator:"required,gte=3"`
    Email string `domainValidator:"required"`
}

func (dto *CreateUserDTO) Validate() (map[string]domain.validateResult, error) {
    return domain.ValidatorInstance().Validate(dto)
}
```

## 3. Main Components

### 3.2. Entities

Entities represent domain objects with identity:

```go
type Product struct {
    *domain.Entity
    Name string
}

func NewProduct(id, name string) *Product {
    return &Product{
        Entity: domain.NewEntity(id),
        Name:   name,
    }
}
```

### 3.3. Aggregate Roots and Domain Events

Aggregates coordinate changes across multiple entities and register domain events:

```go
type UserAggregate struct {
    *domain.AggregateRoot
    Name string
}

func NewUserAggregate(id, name string) *UserAggregate {
    return &UserAggregate{
        AggregateRoot: domain.NewAggregateRoot(id),
        Name:          name,
    }
}
```

To register an event:

```go
type UserCreatedEvent struct {
    id   string
    when time.Time
}

func (e *UserCreatedEvent) Payload() any       { return e }
func (e *UserCreatedEvent) OcurredOn() time.Time { return e.when }
func (e *UserCreatedEvent) Uuid() string       { return e.id }

func (u *UserAggregate) Create() error {
    evt := &UserCreatedEvent{
        id:   u.Uuid(),
        when: time.Now(),
    }
    return u.AddDomainEvent(evt)
}
```

### 3.4. Validation

Validation is driven by struct tags on fields using four built-in validators:

- `required`: value must not be zero (empty string, zero number, nil pointer, etc.)
- `gte=n`: length must be greater than or equal to `n`
- `lte=n`: length must be less than or equal to `n`
- `len=n`: length must be **different** from `n` (intentionally inverted rule)

Example:

```go
type RegisterUserDTO struct {
    Name     string `domainValidator:"required,gte=3"`
    Password string `domainValidator:"required,gte=8"`
}

func (dto *RegisterUserDTO) Validate() (map[string]domain.validateResult, error) {
    return domain.ValidatorInstance().Validate(dto)
}
```

### 3.5. Domain Errors

`errors.go` defines a hierarchy of domain errors:

- `NotFoundError`
- `InternalError`
- `ValidationError`
- `AlreadyExistsError`
- `DependencyError`
- `InvalidDataError`

Example:

```go
func (repo *UserRepository) FindByID(ctx context.Context, id string) (*UserAggregate, error) {
    user, err := repo.load(ctx, id)
    if err == sql.ErrNoRows {
        return nil, domain.NewNotFoundError("user not found")
    }
    if err != nil {
        return nil, domain.NewDependencyError("db failure").SetPreviousError(err)
    }
    return user, nil
}
```

---

From the domain perspective, nothing changes: your handlers still operate on `ddgo`
aggregates and domain events, regardless of whether Kafka or RabbitMQ is used underneath.

