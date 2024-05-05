package config

import (
	"github.com/hex-api-go/internal/user/application/usecase"
	"github.com/hex-api-go/internal/user/domain/contracts"
	"github.com/hex-api-go/internal/user/infrastructure/database"
)

type dependencies struct {
	repository contracts.UserRepository
}

type UserModule struct {
	CreateUserUseCase *usecase.CreateUserUseCase
}

func Bootstrap() *UserModule {
	dependencies := makeDependencies()
	return &UserModule{
		CreateUserUseCase: usecase.NewCreateUserUseCase(dependencies.repository),
	}
}

func makeDependencies() *dependencies {
	return &dependencies{
		repository: database.NewUserRepository(),
	}
}
