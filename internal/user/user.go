package user

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
	"github.com/hex-api-go/pkg/core/infrastructure/rabbitmq"
	gomes "github.com/jeffersonbrasilino/gomes"
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
	gomes.AddActionHandler(createuser.NewComandHandler(userModuleInstance.repository))
	//gomes.AddActionHandler(getuser.NewQueryHandler(nil))
}

func registerPublisher() {
	/* gomes.AddChannelConnection(
		kafka.NewConnection("defaultConKafka", []string{"kafka:9092"}),
	)

	a := kafka.NewPublisherChannelAdapterBuilder(
		"defaultConKafka",
		"gomes.topic",
	)
	//a.WithReplyChannelName("test_response_channel")
	gomes.AddPublisherChannel(a) */

	gomes.AddChannelConnection(
		rabbitmq.NewConnection("rabbit-test", "admin:admin@rabbitmq:5672"),
	)
	pubChan := rabbitmq.NewPublisherChannelAdapterBuilder("rabbit-test", "gomes-exchange").
		WithChannelType(rabbitmq.ProducerQueue)
		//WithExchangeType(rabbitmq.ExchangeFanout).
		//WithExchangeRoutingKeys("fila-1")
	gomes.AddPublisherChannel(pubChan)
}
