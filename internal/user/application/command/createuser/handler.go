package createuser

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jeffersonbrasilino/gomes/otel"
	"github.com/jeffersonbrasilino/hex-api-go/internal/user/domain"
	"github.com/jeffersonbrasilino/hex-api-go/internal/user/domain/contract"
)

type Handler struct {
	repository    contract.UserRepository
	tracer        otel.OtelTrace
	messageHeader map[string]string
}

func NewComandHandler(repository contract.UserRepository) *Handler {
	return &Handler{
		repository: repository,
	}
}

func (c *Handler) Handle(ctx context.Context, data *Command) (any, error) {
	user, errAg := c.makeAggregate(data)
	fmt.Println("data", user, errAg)
	if errAg != nil {
		return nil, errAg
	}

	err := c.repository.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	return "okok", nil
}

func (c *Handler) makeAggregate(data *Command) (*domain.User, error) {
	return domain.NewBuilder().
		WithUuId(uuid.NewString()).
		WithPassword(data.Password).
		WithUsername(data.Username).
		WithPerson(&domain.WithPersonProps{
			Person: &domain.PersonProps{
				UuId:      uuid.NewString(),
				Name:      data.PersonName,
				BirthDate: data.BirthDate,
			},
			Document: &domain.DocumentProps{
				Value: data.Document,
			},
			Contacts: []*domain.ContactProps{
				{
					UuId:        uuid.NewString(),
					ContactType: "email",
					Description: data.Email,
				},
			},
		}).
		Build()
}	
