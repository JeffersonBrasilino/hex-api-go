package createuser

import (
	"github.com/hex-api-go/internal/user/domain/contract"
)

type CommandHandler struct {
	repository contract.UserRepository
}

type ResultCm struct {
	Result any
}

type ParentResultCm struct {
	Parent ResultCm
}

func NewComandHandler(repository contract.UserRepository) *CommandHandler {
	return &CommandHandler{repository}
}

func (c *CommandHandler) Handle(data *Command) (*ResultCm, error) {
	/* return domain.NewUser("new user", "new Password") */
	//fmt.Println("create user > handle CALLED ", data)
	return &ResultCm{"MENSAGEM PROCESSADA COM SUCESSO"}, nil
}
