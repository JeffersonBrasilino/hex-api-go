package database

import (
	"fmt"

	"github.com/hex-api-go/internal/user/domain"
)

type UserRepository struct{}

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}
func (r *UserRepository) Create( aggregate * domain.User) {
	fmt.Println("repo create user")
}
