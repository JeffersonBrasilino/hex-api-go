// User registration command handler.
//
// Intent: orchestrate the creation of a new user, building the domain aggregate from the
// incoming command and persisting it through the repository contract.
// Objective: validate uniqueness by document, build the aggregate via the domain builder, hash
// the raw password before persistence, and delegate creation to the repository.
package createuser

import (
	"context"

	"github.com/google/uuid"
	"github.com/jeffersonbrasilino/ddgo"
	"github.com/jeffersonbrasilino/gomes/otel"
	"github.com/jeffersonbrasilino/hex-api-go/internal/user/domain"
	"github.com/jeffersonbrasilino/hex-api-go/internal/user/domain/contract"
)

// Handler orchestrates the createUser use case: it validates document uniqueness, builds the
// User aggregate, hashes the raw password, and persists the aggregate through the repository.
type Handler struct {
	repository     contract.UserRepository
	passwordHasher contract.PasswordHasher
	tracer         otel.OtelTrace
}

// NewComandHandler creates a Handler for the createUser command.
//
// Intent: wire the handler with the domain contracts it needs to execute the use case.
// Parameters:
//   - repository: persistence contract used to check document uniqueness and create the user.
//   - passwordHasher: contract used to hash the raw password before persistence.
//
// Returns: a ready-to-use *Handler.
func NewComandHandler(
	repository contract.UserRepository,
	passwordHasher contract.PasswordHasher,
) *Handler {
	return &Handler{
		repository:     repository,
		passwordHasher: passwordHasher,
	}
}

// Handle executes the createUser use case.
//
// Intent: register a new user, ensuring the document is not already in use and the password is
// stored hashed rather than in plain text.
// Parameters:
//   - ctx: request-scoped context propagated to the repository and hasher calls.
//   - data: the createUser Command DTO with the user's registration data.
//
// Returns: a success payload on completion, or an error — ddgo.AlreadyExistsError when the
// document is already registered, or any error surfaced by the aggregate builder, the password
// hasher, or the repository.
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

	hash, errHash := c.passwordHasher.Hash(data.Password)
	if errHash != nil {
		return nil, errHash
	}
	user.SetPassword(domain.NewPasswordFromHash(hash))

	err := c.repository.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	return "okok", nil
}

// makeAggregate builds a User aggregate from the createUser Command DTO using the domain
// builder, encapsulating field mapping and optional contact assembly.
//
// Parameters:
//   - data: the createUser Command DTO to map into the aggregate.
//
// Returns: the built *domain.User, or an error if a domain invariant is violated during build.
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
