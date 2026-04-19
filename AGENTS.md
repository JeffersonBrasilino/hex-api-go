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

Always follow the structure and conventions defined in `docs/ARCHITECTURE.md` when:

- creating new files or folders
- organizing code
- naming files, folders, functions, types, etc.
- implementing new features
- creating new modules
- creating new module actions

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

### Principal dependencies documentation

- Gin: https://gin-gonic.com/docs/introduction/
- Gorm: https://gorm.io/docs/
- Gomes: https://github.com/JeffersonBrasilino/gomes
- Ddgo: https://github.com/JeffersonBrasilino/ddgo

### Git conventions and commit messages

Always use conventional commits when:
- creating new commits
- creating new branches
- creating new tags
- creating new pull requests