package createuser

import (
	"github.com/hex-api-go/internal/user/domain"
	"github.com/hex-api-go/internal/user/domain/contract"
)

type CommandHandler struct {
	repository contract.UserRepository
}
type Response struct {
	TestReturn string
}

func NewComandHandler(repository contract.UserRepository) *CommandHandler {
	return &CommandHandler{repository}
}

func (c *CommandHandler) Handle(data *Command) (any, error) {
	return domain.NewUser("new user", "new Password")
}
