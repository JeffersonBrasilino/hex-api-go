package createuser

import (
	"context"
	"fmt"

	"github.com/hex-api-go/internal/user/domain/contract"
	"github.com/jeffersonbrasilino/gomes/otel"
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
		tracer:     otel.InitTrace("command-handler"),
	}
}

func (c *CommandHandler) Handle(ctx context.Context, data *Command) (string, error) {
	fmt.Println("HEADER ACESSOR >>>>>>>", c.messageHeader)
	ctx, span := c.tracer.Start(
		ctx,
		"Handle Command",
		otel.WithMessagingSystemType(otel.MessageSystemTypeInternal),
		otel.WithSpanOperation(otel.SpanOperationProcess),
		otel.WithSpanKind(otel.SpanKindInternal),
	)
	defer span.End()
	return "deu tudo certo", nil
}
