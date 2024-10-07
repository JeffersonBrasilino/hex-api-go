package createuser

import (
	"fmt"

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

func (c *CommandHandler) Handle(data any) (any, error) {
	/* return domain.NewUser("new user", "new Password") */
	fmt.Println("create user > handle CALLED ", data)
	return "hueheuheuh", nil
}

func (c *CommandHandler) Trigger() any {
	return &Command{}
}
