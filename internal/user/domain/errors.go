// Custom domain errors for the user module.
//
// Intent: cover HTTP outcomes not represented by ddgo's fixed error taxonomy
// (401 for invalid/expired sessions, 429 for lockout), keeping the same
// struct + Error() + constructor shape as ddgo's own errors.
package domain

type (
	// InvalidSessionError indicates a session is invalid or expired.
	// It maps to HTTP 401 Unauthorized.
	InvalidSessionError struct {
		message string
	}
	// UserBlockedError indicates a user is temporarily blocked after too
	// many failed login attempts. It maps to HTTP 429 Too Many Requests.
	UserBlockedError struct {
		message string
	}
)

// NewInvalidSessionError creates an error for an invalid or expired session.
//
// Parameters: message — human-readable description.
// Returns: *InvalidSessionError implementing error.
func NewInvalidSessionError(message string) *InvalidSessionError {
	return &InvalidSessionError{message: message}
}

// Error implements the error interface.
func (e *InvalidSessionError) Error() string {
	return e.message
}

// NewUserBlockedError creates an error for a temporarily blocked user.
//
// Parameters: message — human-readable description.
// Returns: *UserBlockedError implementing error.
func NewUserBlockedError(message string) *UserBlockedError {
	return &UserBlockedError{message: message}
}

// Error implements the error interface.
func (e *UserBlockedError) Error() string {
	return e.message
}
