// User module registration.
//
// Intent: compose the user module's dependencies (repositories, adapters, command handlers, and
// HTTP routes) and expose the single entry point (NewUserModule) consumed by main.go.
// Objective: wire the createuser, login, logout, and revokeSession use cases end-to-end — from
// their infrastructure adapters (GORM repository, bcrypt password hasher, JWT token generator,
// Redis-backed login attempt/session stores) through the command bus to their HTTP routes.
package user

import (
	"context"
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/jeffersonbrasilino/gomes"
	_ "github.com/jeffersonbrasilino/gomes/channel/kafka"
	"github.com/jeffersonbrasilino/hex-api-go/internal/user/application/command/createuser"
	"github.com/jeffersonbrasilino/hex-api-go/internal/user/application/command/login"
	"github.com/jeffersonbrasilino/hex-api-go/internal/user/application/command/revokesession"
	"github.com/jeffersonbrasilino/hex-api-go/internal/user/domain/contract"
	"github.com/jeffersonbrasilino/hex-api-go/internal/user/infrastructure/database"
	"github.com/jeffersonbrasilino/hex-api-go/internal/user/infrastructure/http"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// userModule holds the injected infrastructure dependencies (HTTP engine, database connection,
// Redis client, JWT signing secret) and the instantiated contracts shared across the module's
// command handlers.
type userModule struct {
	httpLib         *gin.Engine
	db              *gorm.DB
	redisClient     *redis.Client
	jwtSecret       string
	repository      contract.UserRepository
	dataSource      contract.UserDataSource
	loginRepository contract.LoginRepository
	passwordHasher  contract.PasswordHasher
	tokenGenerator  contract.TokenGenerator
	sessionStore    *database.RedisAdapter
}

// NewUserModule is the constructor called by main.go.
//
// Intent: wire the user module with the infrastructure it needs to construct its adapters —
// the shared Redis client (backing login attempts and session storage) and the JWT signing
// secret (backing access/refresh token issuance). Actual construction of the Redis client and
// resolution of the JWT secret from the environment happen at the call site, outside this
// module.
func NewUserModule(httpLib *gin.Engine, db *gorm.DB, redisClient *redis.Client, jwtSecret string) *userModule {
	return &userModule{
		httpLib:     httpLib,
		db:          db,
		redisClient: redisClient,
		jwtSecret:   jwtSecret,
	}
}

// Register initializes the module's internals and registers its actions/routes.
func (u *userModule) Register(ctx context.Context) error {
	gormRepository := database.NewGormUserRepository(u.db)
	u.repository = gormRepository
	u.loginRepository = gormRepository
	u.passwordHasher = database.NewBcryptAdapter()
	u.tokenGenerator = database.NewJwtAdapter()
	u.sessionStore = database.NewRedisAdapter(u.redisClient)

	u.registerActions()
	u.WithHttpProtocol()
	return nil
}

// WithHttpProtocol defines HTTP routes specific to this module.
func (u *userModule) WithHttpProtocol() *userModule {
	router := u.httpLib.Group("/users")
	http.CreateUserHandler(router)
	http.LoginHandler(router)
	http.RevokeSessionHandler(router)
	slog.Info("User module started with http", "prefix", "/users")
	return u
}

// registerActions maps CQRS commands to their respective handlers.
func (u *userModule) registerActions() {
	gomes.AddActionHandler(createuser.NewComandHandler(u.repository, u.passwordHasher))
	gomes.AddActionHandler(login.NewCommandHandler(
		u.loginRepository,
		u.passwordHasher,
		u.tokenGenerator,
		u.sessionStore,
		u.sessionStore,
	))
	gomes.AddActionHandler(revokesession.NewCommandHandler(u.sessionStore))
}
