package usecase

import (
	"github.com/hex-api-go/internal/user/domain"
	"github.com/hex-api-go/internal/user/domain/contract"
	"github.com/hex-api-go/pkg/core"
)

type CreateUserUseCase struct {
	repository contract.UserRepository
}

func NewCreateUserUseCase(repository contract.UserRepository) *CreateUserUseCase {
	return &CreateUserUseCase{
		repository: repository,
	}
}

func (u *CreateUserUseCase) CreateUser(data interface{}) *core.Result {
	user := domain.NewUser("new user", "new Password")
	return core.ResultSuccess(user)
}
