#### Module Register Pattern

Below is the boilerplate standard for implementing the module registration file. It demonstrates dependency injection layout and registration bootstrapping on a single file.

```go
package [module-name]

import (
	"context"
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/jeffersonbrasilino/gomes"
	// Import necessary internal components
	"github.com/jeffersonbrasilino/hex-api-go/internal/[module-name]/application/command/[actionname]"
	"github.com/jeffersonbrasilino/hex-api-go/internal/[module-name]/domain/contract"
	"github.com/jeffersonbrasilino/hex-api-go/internal/[module-name]/infrastructure/database"
	"github.com/jeffersonbrasilino/hex-api-go/internal/[module-name]/infrastructure/http"
	"gorm.io/gorm"
)

// [moduleName]Module internal struct holding injected dependencies and instantiated contracts.
type [moduleName]Module struct {
	httpLib    *gin.Engine
	db         *gorm.DB
	repository contract.[AggregateRoot]Repository
}

// New[ModuleName]Module is the constructor called by main.go.
func New[ModuleName]Module(httpLib *gin.Engine, db *gorm.DB) *[moduleName]Module {
	return &[moduleName]Module{
		httpLib: httpLib,
		db:      db,
	}
}

// Register initializes the module's internals and registers its actions/routes.
func (m *[moduleName]Module) Register(ctx context.Context) error {
	// 1. Initialize Infrastructural Dependencies (e.g. Repositories)
	m.repository = database.NewGorm[ModuleName]Repository(m.db)

	// 2. Register CQRS Handlers
	m.registerActions()

	// 3. Setup HTTP Routing
	m.WithHttpProtocol()

	return nil
}

// WithHttpProtocol defines HTTP routes specific to this module.
func (m *[moduleName]Module) WithHttpProtocol() *[moduleName]Module {
	router := m.httpLib.Group("/[module-route-prefix]") // e.g. "/users"

	// Register handlers on the group
	http.Create[ModuleName]Handler(router)

	slog.Info("[ModuleName] module started with http", "prefix", "/[module-route-prefix]")
	return m
}

// registerActions maps CQRS commands/queries to their respective handlers.
func (m *[moduleName]Module) registerActions() {
	gomes.AddActionHandler([actionname].NewCommandHandler(m.repository))
}
```

Implementation logic reference: see -> `internal/user/user.go`
