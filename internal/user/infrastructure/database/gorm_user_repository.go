package database

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/jeffersonbrasilino/ddgo"
	"github.com/jeffersonbrasilino/hex-api-go/internal/user/domain"
	"gorm.io/gorm"
)

type GormUserRepository struct {
	db *gorm.DB
}

func NewGormUserRepository(db *gorm.DB) *GormUserRepository {

	if os.Getenv("GORM_AUTO_MIGRATE") == "1" {
		db.SetupJoinTable(&Users{}, "UserGroups", &UserGroupUser{})
		err := db.AutoMigrate(&Users{}, &Person{}, &UsersGroups{}, &PersonContacts{}, &PersonContactsType{}, &UserGroupsPermissions{}, &UsersDevice{})
		if err != nil {
			slog.Error("[GormUserRepository]", "error", err)
		}
	}

	if os.Getenv("APP_ENV") == "local" {
		db = db.Debug()
	}

	return &GormUserRepository{db: db}
}

func (r *GormUserRepository) Create(ctx context.Context, user *domain.User) error {
	tx := r.db.Begin()
	entity := toDatabase(user)
	result := gorm.WithResult()
	err := gorm.G[Users](tx, result).Create(ctx, entity)
	if err != nil {
		tx.Rollback()
		return ddgo.NewInternalError(fmt.Sprintf("Error to create user: %s", err.Error()))
	}

	fmt.Println("[REPOSITORY RESULT]", result)

	return tx.Commit().Error
}
