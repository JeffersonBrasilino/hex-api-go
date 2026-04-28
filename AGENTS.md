### Project Context

This project applies hexagonal architecture, Domain-Driven Design (DDD), and Event-driven architecture concepts with Golang. It uses the following dependencies: Gin (HTTP server), GORM (database), Gomes (CQRS and event-sourcing), Ddgo (DDD). It utilizes a modular monolith approach for context separation.

### Project commands

- `make deps` - install dependencies (vendor)
- `make start-dev` - start dev environment
- `make test` - run tests
- `make coverage-terminal` - run tests with terminal coverage
- `make coverage-html` - run tests with html coverage
- `make pprof-goroutine` - run pprof with goroutine
- `make pprof-cpu` - run pprof with cpu
- `make pprof-heap` - run pprof with heap

### Project structure and architecture

This project follows a strict modular monolith structure where each module under `internal/`
has three layers: domain, application, and infrastructure. All modules must follow the same
structure, naming conventions, and implementation patterns. When working with module
composition, layer boundaries, or component implementation, consult the appropriate skill
for detailed guidelines.

```
[project-name]/
├── cmd/                                # entry points of the application
│   └── api/
│       └── main.go
├── internal/                           # private code of the application
│   ├── [module-name]/                  # module of the application
│   │   ├── domain/                     # business rules and contracts
│   │   ├── application/                # CQRS command/query handlers
│   │   └── infrastructure/             # adapters (database, http, etc)
│   └── [module-name]/[module-name].go  # module registration file
├── pkg/                                # public/shared code
├── docs/                               # project documentation
├── .agents/                            # AI agent configuration
│   └── skills/                         # agent skills (knowledge + actions)
└── vendor/                             # vendored dependencies
```

### Principal dependencies

- **Gin** (HTTP server): https://gin-gonic.com/docs/
- **GORM** (database ORM): https://gorm.io/docs/
- **Gomes** (CQRS/Event bus): https://github.com/JeffersonBrasilino/gomes
- **Ddgo** (DDD primitives): https://github.com/JeffersonBrasilino/ddgo

### Safety and permissions

Allowed without prompt:
- read files, list files
- go fmt ./..., go vet ./...
- make test

Ask first:
- package installs
- deleting files, chmod
- running full build or end to end suites

### After codebase changes

After codebase changes, execute:

- `go fmt ./...` - lint Golang code
- `go vet ./...` - check for compilation errors
- `make test` - run all tests
- `make coverage-terminal` - run all tests with code coverage

### Git conventions and commit messages

Always use conventional commits when:
- creating new commits
- creating new branches
- creating new tags
- creating new pull requests