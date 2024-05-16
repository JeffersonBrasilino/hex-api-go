package database

import (
	"fmt"
	"time"

	"github.com/hex-api-go/internal/user/domain"
)

type UserRepository struct{}

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}
func (r *UserRepository) Create(aggregate *domain.User) {
	time.Sleep(time.Second * 3)
	fmt.Println("repo create user")
}
