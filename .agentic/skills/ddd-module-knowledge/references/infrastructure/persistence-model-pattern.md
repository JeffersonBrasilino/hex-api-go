#### Persistence Model pattern (GORM)

Persistence models are infrastructure-specific structs that map directly to database tables.
They must be completely separate from domain entities to maintain layer isolation.

The persistence model must:

- embed `gorm.Model` for standard fields (`ID`, `CreatedAt`, `UpdatedAt`, `DeletedAt`).
- use `gorm` struct tags for column mapping, constraints, and relationships.
- define a `TableName()` method with the pattern `[project-name].[table_name]`.
- use plural names for structs that represent database tables ex: `Users`, `PersonContacts`.
- define foreign keys explicitly in struct tags.
- use `many2many` tag with explicit `joinForeignKey` and `joinReferences` for many-to-many relationships.
- be defined in the same file `gorm_model.go` for all models of the module.

Boilerplate Example:

```go
package database

import (
	"time"

	"gorm.io/gorm"
)

// main entity model
type [ModulePlural] struct {
	gorm.Model
	Field1     string           `gorm:"column:field1;not null"`
	Field2     string           `gorm:"column:field2;not null"`
	ChildId    uint             `gorm:"column:child_id;not null"`
	Child      [ChildModel]
	Relations  [][RelationModel] `gorm:"many2many:[join_table];joinForeignKey:[fk_column];joinReferences:[ref_column]"`
}

// child entity model
type [ChildModel] struct {
	gorm.Model
	Name      string            `gorm:"column:name;not null"`
	Parents   [][ModulePlural]  `gorm:"foreignKey:[ParentFK]"`
}

// join table model (for many-to-many with extra fields)
type [JoinModel] struct {
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     gorm.DeletedAt `gorm:"index"`
	ExtraField    bool           `gorm:"column:extra_field;not null; default:false"`
	LeftID        uint           `gorm:"column:left_id;primaryKey"`
	RightID       uint           `gorm:"column:right_id;primaryKey"`
}

// TableName methods - pattern: [project-name].[table_name]
func ([ModulePlural]) TableName() string {
	return "hex-api-go.[table_name]"
}

func ([ChildModel]) TableName() string {
	return "hex-api-go.[child_table_name]"
}

func ([JoinModel]) TableName() string {
	return "hex-api-go.[join_table_name]"
}
```

Implementation example: see -> `internal/user/infrastructure/database/gorm_model.go`
