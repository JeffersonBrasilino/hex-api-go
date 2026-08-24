// Login attempt store contract.
//
// Intent: define control of failed login attempts to enforce a lockout
// window after repeated failures.
// Objective: allow login to register failures, check lockout status, and
// reset the counter on success.
package contract

import "context"

// LoginAttemptStore defines control of failed login attempts and lockout
// state for the login user action.
type LoginAttemptStore interface {
	// RegisterFailure records a failed login attempt for the given username.
	RegisterFailure(ctx context.Context, username string) error
	// IsBlocked reports whether the given username is currently locked out.
	IsBlocked(ctx context.Context, username string) (bool, error)
	// Reset clears the failed attempt count for the given username.
	Reset(ctx context.Context, username string) error
}
