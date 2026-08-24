// Login session store contract.
//
// Intent: define creation of a login session record.
// Objective: allow login to persist a new session keyed by session id, with
// an expiration, scoped exclusively to session creation.
package contract

import (
	"context"
)

// LoginSessionStore defines creation of a login session for the login user
// action.
type LoginSessionStore interface {
	// Save persists a session identified by sessionId and refreshToken, expiring after ttl.
	Save(ctx context.Context, sessionId string, refreshToken string) error
}
