package user

import (
	"context"
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
	messagesystem "github.com/hex-api-go/pkg/core/infrastructure/message_system"
	"github.com/hex-api-go/pkg/core/infrastructure/message_system/channel/kafka"
)

var userModuleInstance *userModule

type userModule struct {
	repository contract.UserRepository
	dataSource contract.UserDataSource
}

func Bootstrap() *userModule {

	if userModuleInstance != nil {
		return userModuleInstance
	}

	userModuleInstance = &userModule{
		repository: database.NewUserRepository(),
		dataSource: facade.NewUserFacade(makeAclGateways()),
	}
	registerActions()

	return userModuleInstance
}

func makeAclGateways() map[string]aclcontract.PersonGateway {
	return map[string]aclcontract.PersonGateway{
		"gatewayA": gateway.NewJsonPlaceholderGateway(),
		"gatewayB": gateway.NewRandonUserMeGateway(),
	}
}

func (u *userModule) WithHttpProtocol(ctx context.Context, httpLib *fiber.App) *userModule {

	registerPublisher()

	router := httpLib.Group("/users")
	http.CreateUser(ctx, router)
	fmt.Println("User module started with http. Prefix: /users")
	return u
}

func registerActions() {
	messagesystem.AddActionHandler(createuser.NewComandHandler(userModuleInstance.repository))
	messagesystem.AddActionHandler(getuser.NewQueryHandler(nil))
}

func registerPublisher() {
	messagesystem.AddChannelConnection(
		kafka.NewConnection("defaultConKafka", []string{"kafka:9092"}),
	)

	a := kafka.NewPublisherChannelAdapterBuilder(
		"defaultConKafka",
		"message_system.topic",
	)
	a.WithReplyChannelName("test_response_channel")
	
	messagesystem.AddPublisherChannel(a)
}
