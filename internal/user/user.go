package user

import (
	"context"
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/jeffersonbrasilino/gomes"
	_ "github.com/jeffersonbrasilino/gomes/channel/kafka"
	"github.com/jeffersonbrasilino/hex-api-go/internal/user/application/command/createuser"
	"github.com/jeffersonbrasilino/hex-api-go/internal/user/domain/contract"
	aclcontract "github.com/jeffersonbrasilino/hex-api-go/internal/user/infrastructure/acl/contract"
	"github.com/jeffersonbrasilino/hex-api-go/internal/user/infrastructure/acl/facade"
	"github.com/jeffersonbrasilino/hex-api-go/internal/user/infrastructure/acl/gateway"
	"github.com/jeffersonbrasilino/hex-api-go/internal/user/infrastructure/database"
	"github.com/jeffersonbrasilino/hex-api-go/internal/user/infrastructure/http"
)

type userModule struct {
	httpLib    *gin.Engine
	db         any
	repository contract.UserRepository
	dataSource contract.UserDataSource
}

func NewUserModule(httpLib *gin.Engine, db any) *userModule {
	return &userModule{
		httpLib: httpLib,
		db:      db,
	}
}

func (u *userModule) Register(ctx context.Context) error {
	u.repository = database.NewUserRepository()
	u.dataSource = facade.NewUserFacade(
		map[string]aclcontract.PersonGateway{
			"gatewayA": gateway.NewJsonPlaceholderGateway(),
			"gatewayB": gateway.NewRandonUserMeGateway(),
		},
	)

	u.registerActions()
	u.WithHttpProtocol()
	return nil
}

func (u *userModule) WithHttpProtocol() *userModule {
	router := u.httpLib.Group("/users")
	http.CreateUserHandler(router)
	slog.Info("User module started with http", "prefix", "/users")
	return u
}

func (u *userModule) registerActions() {
	gomes.AddActionHandler(createuser.NewComandHandler(u.repository))
}
