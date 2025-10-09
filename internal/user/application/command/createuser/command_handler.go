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
	return &CommandHandler{repository}
}

func (c *CommandHandler) Handle(ctx context.Context, data *Command) (*ResultCm, error) {

	fmt.Println("iniciando regra de negocio", data)
	//time.Sleep(time.Second * 3)

	//return &ResultCm{"MENSAGEM PROCESSADA COM SUCESSO"}, fmt.Errorf("deu ruim ao processar a mensagem")

	return &ResultCm{"MENSAGEM PROCESSADA COM SUCESSO"}, nil
}
