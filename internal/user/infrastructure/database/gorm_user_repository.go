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
		//db.SetupJoinTable(&Users{}, "UserGroups", &UserGroupUser{})
		err := db.AutoMigrate(&Users{}, &Person{}, &UsersGroups{}, &PersonContacts{}, &PersonContactsType{}, &UserGroupsPermissions{}, &UsersDevice{}, &UserGroupUser{})
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

	contacts := []PersonContacts{}
	for _, c := range user.Person().Contacts() {
		ct, err := gorm.G[PersonContactsType](r.db).
			Select("id").
			Where("uuid = ?", c.ContactType().Uuid()).
			First(ctx)
		if err != nil {
			tx.Rollback()
			return ddgo.NewInternalError(fmt.Sprintf("Error to create user: %s", err.Error()))
		}
		contacts = append(contacts, PersonContacts{
			Uuid:          c.Uuid(),
			ContactTypeId: ct.ID,
			Contact:       c.Description(),
			Main:          true,
			Status:        1,
		})
	}
	entity.Person.Contacts = contacts

	userGroups := []UserGroupUser{}
	for _, ug := range user.UserGroups() {
		ugData, err := gorm.G[UsersGroups](r.db).Select("id").
			Where("uuid = ?", ug.Uuid()).
			First(ctx)
		if err != nil {
			tx.Rollback()
			return ddgo.NewInternalError(fmt.Sprintf("Error to create user: %s", err.Error()))
		}
		userGroups = append(userGroups, UserGroupUser{
			Uuid:        ug.Uuid(),
			Status:      1,
			Main:        true,
			UserGroupId: ugData.ID,
		})
	}
	entity.UserGroupsUsers = userGroups

	result := gorm.WithResult()
	err := gorm.G[Users](tx, result).Create(ctx, entity)
	if err != nil {
		tx.Rollback()
		return ddgo.NewInternalError(fmt.Sprintf("Error to create user: %s", err.Error()))
	}

	fmt.Println("[REPOSITORY RESULT]")

	return tx.Commit().Error
}

func (r *GormUserRepository) ExistsByDocument(ctx context.Context, document string) (bool, error) {
	users, err := gorm.G[Person](r.db, nil).
		Where("document = ?", document).
		Count(ctx, "document")
	if err != nil {
		return false, ddgo.NewInternalError(
			fmt.Sprintf("Error to check if exists user by document: %s", err.Error()),
		)
	}
	return users > 0, nil
}
