// Revoke session store contract.
//
// Intent: define removal of an existing session.
// Objective: allow logout and administrative session revocation to share a
// single deletion contract, since both perform the same domain action.
package contract

import "context"

// RevokeSessionStore defines the shared session revocation action used by
// both logout and administrative session revocation.
type RevokeSessionStore interface {
	// Delete removes the session identified by sessionId. It returns a
	// "not found" error when the session does not exist.
	Delete(ctx context.Context, sessionId string) error
}
