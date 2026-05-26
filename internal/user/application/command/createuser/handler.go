package createuser

import (
	"context"

	"github.com/google/uuid"
	"github.com/jeffersonbrasilino/ddgo"
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
	
	exists, errExists := c.repository.ExistsByDocument(ctx, data.Document)
	if errExists != nil {
		return nil, errExists
	}
	if exists {
		return nil, ddgo.NewAlreadyExistsError("User already exists")
	}

	user, errAg := c.makeAggregate(data)
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

	contactData := []*domain.WithContactProps{}
	if data.Email != "" {
		contactData = append(contactData, &domain.WithContactProps{
			UuId:        uuid.NewString(),
			Description: data.Email,
			ContactType: "f70e57f1-244a-4ef7-ab27-05f5adc777d7",
		})
	}

	return domain.NewBuilder().
		WithUuId(uuid.NewString()).
		WithPassword(data.Password).
		WithUsername(data.Username).
		WithPerson(&domain.WithPersonProps{
			UuId:      uuid.NewString(),
			Name:      data.PersonName,
			BirthDate: data.BirthDate,
			Document:  data.Document,
			Contacts:  contactData,
		}).
		WithUserGroups([]*domain.UserGroupProps{
			{
				UuId: "422eacba-efda-4c0a-af22-cf3b2f92b174",
			},
		}).
		Build()
}
