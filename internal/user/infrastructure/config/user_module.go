package config

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/hex-api-go/internal/user/application/command/createuser"
	"github.com/hex-api-go/internal/user/application/query/getuser"
	"github.com/hex-api-go/internal/user/domain/contract"
	aclcontract "github.com/hex-api-go/internal/user/infrastructure/acl/contract"
	"github.com/hex-api-go/internal/user/infrastructure/acl/facade"
	"github.com/hex-api-go/internal/user/infrastructure/acl/gateway"
	"github.com/hex-api-go/internal/user/infrastructure/database"
	"github.com/hex-api-go/internal/user/infrastructure/http"
	"github.com/hex-api-go/pkg/core/application/cqrs"
)

type dependencies struct {
	repository contract.UserRepository
	dataSource contract.UserDataSource
}

func Bootstrap() {
	dependencies := makeDependencies()
	registerActions(dependencies)
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

func registerActions(dependencies *dependencies) {
	cqrs.RegisterActionHandler(createuser.NewComandHandler(dependencies.repository))
	cqrs.RegisterActionHandler(getuser.NewQueryHandler(dependencies.dataSource))
}

func WithHttpHandlers(fiberApp *fiber.App) {
	http.NewHttpHandlers(fiberApp)
	fmt.Println("User module started with http.")
}
