package createuser

import (
	"context"
	"fmt"

	"github.com/hex-api-go/internal/user/domain/contract"
)

type CommandHandler struct {
	repository contract.UserRepository
}

type ResultCm struct {
	Result any
}

func NewComandHandler(repository contract.UserRepository) *CommandHandler {
	return &CommandHandler{repository: repository}
}

func (c *CommandHandler) Handle(ctx context.Context, data *Command) (string, error) {
	return "deu tudo certo", fmt.Errorf("deu ruim ao processar a mensagem")
}
