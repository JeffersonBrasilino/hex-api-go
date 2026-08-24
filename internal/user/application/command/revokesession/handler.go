// Revoke session command handler.
//
// Intent: orchestrate the administrative termination of any user's active session.
// Objective: remove the session identified by the id informed by the administrator, so it can no
// longer be used to access the system.
package revokesession

import (
	"context"

	"github.com/jeffersonbrasilino/ddgo"
	"github.com/jeffersonbrasilino/hex-api-go/internal/user/domain/contract"
)

// Handler orchestrates the administrative revokeSession use case: it deletes the session
// identified by the informed session id through the revoke session store.
type Handler struct {
	revokeSessionStore contract.RevokeSessionStore
}

// NewCommandHandler creates a Handler for the revokeSession command.
//
// Intent: wire the handler with the domain contract it needs to execute the use case.
// Parameters:
//   - revokeSessionStore: contract used to delete the session.
//
// Returns: a ready-to-use *Handler.
func NewCommandHandler(revokeSessionStore contract.RevokeSessionStore) *Handler {
	return &Handler{
		revokeSessionStore: revokeSessionStore,
	}
}

// Handle executes the revokeSession use case.
//
// Intent: end the session identified by the administrator-provided session id, so it can no
// longer be used to access the system.
// Parameters:
//   - ctx: request-scoped context propagated to the contract.
//   - data: the revokeSession Command DTO with the session id to revoke.
//
// Returns: a success payload on completion, or a ddgo.NotFoundError (404) when the session does
// not exist.
func (c *Handler) Handle(ctx context.Context, data *Command) (any, error) {
	if errDelete := c.revokeSessionStore.Delete(ctx, data.SessionId); errDelete != nil {
		return nil, ddgo.NewNotFoundError("Sessão não encontrada")
	}

	return "session revoked", nil
}
