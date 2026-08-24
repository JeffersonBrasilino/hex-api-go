// Login command handler.
//
// Intent: orchestrate user authentication, enforcing lockout after repeated failed attempts and
// issuing a session with access/refresh tokens on success.
// Objective: check the lockout state, verify credentials, register or reset failed attempts, and
// on success create a session and issue tokens for the authenticated user.
package login

import (
	"context"

	"github.com/google/uuid"
	"github.com/jeffersonbrasilino/ddgo"
	"github.com/jeffersonbrasilino/hex-api-go/internal/user/domain"
	"github.com/jeffersonbrasilino/hex-api-go/internal/user/domain/contract"
)

// Handler orchestrates the login use case: it checks lockout state, verifies credentials against
// the stored user, registers or resets failed attempts, and issues a session with access and
// refresh tokens on success.
type Handler struct {
	loginRepository   contract.LoginRepository
	passwordHasher    contract.PasswordHasher
	tokenGenerator    contract.TokenGenerator
	loginAttemptStore contract.LoginAttemptStore
	loginSessionStore contract.LoginSessionStore
}

// NewCommandHandler creates a Handler for the login command.
//
// Intent: wire the handler with the domain contracts it needs to execute the use case.
// Parameters:
//   - loginRepository: contract used to locate the user by username or document.
//   - passwordHasher: contract used to verify the raw password against the stored hash.
//   - tokenGenerator: contract used to issue access and refresh tokens.
//   - loginAttemptStore: contract used to check, register, and reset failed login attempts.
//   - loginSessionStore: contract used to persist the newly created session.
//
// Returns: a ready-to-use *Handler.
func NewCommandHandler(
	loginRepository contract.LoginRepository,
	passwordHasher contract.PasswordHasher,
	tokenGenerator contract.TokenGenerator,
	loginAttemptStore contract.LoginAttemptStore,
	loginSessionStore contract.LoginSessionStore,
) *Handler {
	return &Handler{
		loginRepository:   loginRepository,
		passwordHasher:    passwordHasher,
		tokenGenerator:    tokenGenerator,
		loginAttemptStore: loginAttemptStore,
		loginSessionStore: loginSessionStore,
	}
}

// Handle executes the login use case.
//
// Intent: authenticate the user identified by the command, enforcing the lockout window after
// repeated failures, and issue a session with access/refresh tokens on success.
// Parameters:
//   - ctx: request-scoped context propagated to the contracts.
//   - data: the login Command DTO with the submitted username and password.
//
// Returns: a map with accessToken, refreshToken, and sessionId on success, or an error —
// domain.UserBlockedError when the user is locked out, ddgo.ValidationError when the credentials
// are invalid, or any error surfaced by the underlying contracts.
func (c *Handler) Handle(ctx context.Context, data *Command) (any, error) {
	blocked, errBlocked := c.loginAttemptStore.IsBlocked(ctx, data.Username)
	if errBlocked != nil {
		return nil, errBlocked
	}
	if blocked {
		return nil, domain.NewUserBlockedError(
			"Usuário temporariamente bloqueado por excesso de tentativas",
		)
	}

	user, errUser := c.loginRepository.FindByUsernameOrDocument(ctx, data.Username)
	if errUser != nil || user == nil {
		if errAttempt := c.loginAttemptStore.RegisterFailure(ctx, data.Username); errAttempt != nil {
			return nil, errAttempt
		}
		return nil, ddgo.NewValidationError("Credenciais inválidas")
	}

	valid, errVerify := c.passwordHasher.Verify(data.Password, user.Password().Value())
	if errVerify != nil {
		return nil, errVerify
	}
	
	if !valid {
		if errAttempt := c.loginAttemptStore.RegisterFailure(ctx, data.Username); errAttempt != nil {
			return nil, errAttempt
		}
		return nil, ddgo.NewValidationError("Credenciais inválidas")
	}

	if errReset := c.loginAttemptStore.Reset(ctx, data.Username); errReset != nil {
		return nil, errReset
	}

	sessionId := uuid.NewString()

	groups := make([]string, 0, len(user.UserGroups()))
	for _, group := range user.UserGroups() {
		groups = append(groups, group.Name())
	}

	accessToken, errAccessToken := c.tokenGenerator.GenerateAccessToken(user.Uuid(), groups)
	if errAccessToken != nil {
		return nil, errAccessToken
	}

	refreshToken, errRefreshToken := c.tokenGenerator.GenerateRefreshToken(sessionId)
	if errRefreshToken != nil {
		return nil, errRefreshToken
	}

	if errSave := c.loginSessionStore.Save(ctx, sessionId, refreshToken); errSave != nil {
		return nil, errSave
	}

	return map[string]string{
		"accessToken":  accessToken,
		"refreshToken": sessionId,
	}, nil
}
