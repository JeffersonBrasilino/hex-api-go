// Login repository contract.
//
// Intent: define read access needed to locate a user during the login flow.
// Objective: allow the login use case to find a user by username or document
// without depending on the broader UserRepository used by other actions.
package contract

import (
	"context"

	"github.com/jeffersonbrasilino/hex-api-go/internal/user/domain"
)

// LoginRepository defines lookup access to a User for the login user action.
type LoginRepository interface {
	// FindByUsernameOrDocument returns the user matching the given username or document.
	FindByUsernameOrDocument(ctx context.Context, identifier string) (*domain.User, error)
}
