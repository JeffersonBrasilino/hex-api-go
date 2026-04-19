#### Repository Implementation pattern

Repository implementations are the concrete adapters that fulfill domain contracts (interfaces).
They encapsulate all persistence logic and translate between domain aggregates and the underlying storage mechanism.

The repository must:
- implement the domain contract interface defined in `domain/contract/`.
- receive the database dependency via constructor injection.
- use transaction management (`Begin`, `Commit`, `Rollback`) for write operations.
- delegate domain ↔ persistence conversion to mapper functions.
- wrap errors with `ddgo` error types.

Boilerplate Example:

```go
package database

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/jeffersonbrasilino/ddgo"
	"github.com/jeffersonbrasilino/hex-api-go/internal/[module-name]/domain"
	"gorm.io/gorm"
)

type Gorm[ModuleName]Repository struct {
	db *gorm.DB
}

func NewGorm[ModuleName]Repository(db *gorm.DB) *Gorm[ModuleName]Repository {

	if os.Getenv("GORM_AUTO_MIGRATE") == "1" {
		// register join tables before migration if needed
		// db.SetupJoinTable(&Model{}, "Relation", &JoinModel{})
		err := db.AutoMigrate(&Model1{}, &Model2{})
		if err != nil {
			slog.Error("[Gorm[ModuleName]Repository]", "error", err)
		}
	}

	if os.Getenv("APP_ENV") == "local" {
		db = db.Debug()
	}

	return &Gorm[ModuleName]Repository{db: db}
}

func (r *Gorm[ModuleName]Repository) Create(ctx context.Context, aggregate *domain.[AggregateRoot]) error {
	tx := r.db.Begin()
	entity := toDatabase(aggregate)
	result := gorm.WithResult()
	err := gorm.G[[PersistenceModel]](tx, result).Create(ctx, entity)
	if err != nil {
		tx.Rollback()
		return ddgo.NewInternalError(fmt.Sprintf("Error to create [module-name]: %s", err.Error()))
	}

	return tx.Commit().Error
}
```
Implementation example: see -> `../../internal/user/infrastructure/database/gorm_user_repository.go`
