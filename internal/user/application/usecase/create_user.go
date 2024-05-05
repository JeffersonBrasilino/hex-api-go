package usecase

import (
	"fmt"

	"github.com/hex-api-go/internal/user/domain"
	"github.com/hex-api-go/internal/user/domain/contracts"
)

type CreateUserUseCase struct {
	repository contracts.UserRepository
}

func NewCreateUserUseCase(repository contracts.UserRepository) *CreateUserUseCase {
	return &CreateUserUseCase{
		repository: repository,
	}
}
func (u *CreateUserUseCase) CreateUser(data interface{}) {
	user := domain.NewUser("new user", "new Password")
	//u.repository.Create(user)
	fmt.Println("#### use case ######")
	fmt.Println(user)
}
