package createuser

import (
	"context"

	"github.com/jeffersonbrasilino/gomes/otel"
	"github.com/jeffersonbrasilino/hex-api-go/internal/user/domain"
	"github.com/jeffersonbrasilino/hex-api-go/internal/user/domain/contract"
)

type CommandHandler struct {
	repository    contract.UserRepository
	tracer        otel.OtelTrace
	messageHeader map[string]string
}
type ResultCm struct {
	Result any
}

func NewComandHandler(repository contract.UserRepository) *CommandHandler {
	return &CommandHandler{
		repository: repository,
	}
}

func (c *CommandHandler) Handle(ctx context.Context, data *Command) (string, error) {
	props := &domain.UserProps{}
	_, err := domain.NewUser(props)

	if err != nil {
		return "", err
	}
	return "okok", nil
}
