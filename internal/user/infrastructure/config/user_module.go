package config

import (
	"github.com/hex-api-go/internal/user/application/usecase"
	"github.com/hex-api-go/internal/user/domain/contract"
	aclcontract "github.com/hex-api-go/internal/user/infrastructure/acl/contract"
	"github.com/hex-api-go/internal/user/infrastructure/acl/facade"
	"github.com/hex-api-go/internal/user/infrastructure/acl/gateway"
	"github.com/hex-api-go/internal/user/infrastructure/database"
)

type dependencies struct {
	repository contract.UserRepository
	dataSource contract.UserDataSource
}

type UserModule struct {
	CreateUserUseCase *usecase.CreateUserUseCase
	GetUserUseCase    *usecase.GetUserUseCase
}

func Bootstrap() *UserModule {
	dependencies := makeDependencies()
	return &UserModule{
		CreateUserUseCase: usecase.NewCreateUserUseCase(dependencies.repository),
		GetUserUseCase:    usecase.NewGetUserUseCase(dependencies.dataSource),
	}
}

func makeDependencies() *dependencies {
	gatewaysAcl := map[string]aclcontract.PersonGateway{
		"gatewayA": gateway.NewJsonPlaceholderGateway(),
		"gatewayB": gateway.NewRandonUserMeGateway(),
	}
	return &dependencies{
		repository: database.NewUserRepository(),
		dataSource: facade.NewUserFacade(gatewaysAcl),
	}
}
