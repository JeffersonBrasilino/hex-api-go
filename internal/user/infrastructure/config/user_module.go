package config

import (
	"context"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/hex-api-go/internal/user/application/command/createuser"
	"github.com/hex-api-go/internal/user/domain/contract"
	aclcontract "github.com/hex-api-go/internal/user/infrastructure/acl/contract"
	"github.com/hex-api-go/internal/user/infrastructure/acl/facade"
	"github.com/hex-api-go/internal/user/infrastructure/acl/gateway"
	"github.com/hex-api-go/internal/user/infrastructure/database"
	"github.com/hex-api-go/internal/user/infrastructure/http"
	"github.com/hex-api-go/pkg/core/infrastructure/message_system/bus"
)

type userModule struct {
	repository contract.UserRepository
	dataSource contract.UserDataSource
}

var userModuleInstance *userModule

func bootstrap() {

	if userModuleInstance != nil {
		return
	}

	userModuleInstance = &userModule{
		repository: database.NewUserRepository(),
		dataSource: facade.NewUserFacade(makeAclGateways()),
	}
	defer registerActions()
}

func StartModuleWithHttpServer(ctx context.Context, fiberApp *fiber.App) {
	bootstrap()
	router := fiberApp.Group("/users")
	http.CreateUser(ctx, router)
	fmt.Println("User module started with http. Prefix: /users")
}

func makeAclGateways() map[string]aclcontract.PersonGateway {
	return map[string]aclcontract.PersonGateway{
		"gatewayA": gateway.NewJsonPlaceholderGateway(),
		"gatewayB": gateway.NewRandonUserMeGateway(),
	}
}

func registerActions() {
	//messagesystem.AddCommandHandler("CreateUser", createuser.NewComandHandler(userModuleInstance.repository))
	//messageSystem.AddCommandHandler("GetUser", getuser.NewQueryHandler(nil))
	//messageSystem.AddQueryHandler("GetUser", getuser.NewQueryHandler(nil))
	//cqrs.RegisterActionHandler(createuser.NewComandHandler(userModuleInstance.repository))
	//cqrs.RegisterActionHandler(getuser.NewQueryHandler(dependencies.dataSource))

	bus.RegisterCommandHandler("CreateUser", createuser.NewComandHandler(userModuleInstance.repository))
}
