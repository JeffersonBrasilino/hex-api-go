// Custom domain errors for consistent handling across the application.
//
// Intent: Standardize error reporting so all layers use known error types
// that can be mapped to HTTP status codes or other responses.
// Objective: Make it easier to detect and handle domain errors (e.g. not
// found, validation, dependency failure) without string matching.
package domain

import "fmt"

type (
	abstractError struct {
		message  string
		previous error
	}
	// NotFoundError indicates a requested resource does not exist.
	NotFoundError struct {
		abstractError
	}
	// InternalError indicates an unexpected internal failure.
	InternalError struct {
		abstractError
	}
	// ValidationError indicates input or state failed validation.
	ValidationError struct {
		abstractError
	}
	// AlreadyExistsError indicates a duplicate or conflicting resource.
	AlreadyExistsError struct {
		abstractError
	}
	// DependencyError indicates a failure in an external dependency.
	DependencyError struct {
		abstractError
	}
	// InvalidDataError indicates the provided data is invalid.
	InvalidDataError struct {
		abstractError
	}
)

// NewDependencyError creates an error for dependency failures.
//
// Parameters: message — human-readable description.
// Returns: *DependencyError implementing error.
func NewDependencyError(message string) *DependencyError {
	return &DependencyError{
		abstractError{
			message: message,
		},
	}
}

// NewNotFoundError creates an error for missing resources.
//
// Parameters: message — human-readable description.
// Returns: *NotFoundError implementing error.
func NewNotFoundError(message string) *NotFoundError {
	return &NotFoundError{
		abstractError{
			message: message,
		},
	}
}

// NewInternalError creates an error for internal failures.
//
// Parameters: message — human-readable description.
// Returns: *InternalError implementing error.
func NewInternalError(message string) *InternalError {
	return &InternalError{
		abstractError{
			message: message,
		},
	}
}

// NewAlreadyExistsError creates an error for duplicate resources.
//
// Parameters: message — human-readable description.
// Returns: *AlreadyExistsError implementing error.
func NewAlreadyExistsError(message string) *AlreadyExistsError {
	return &AlreadyExistsError{
		abstractError{
			message: message,
		},
	}
}

// NewInvalidDataError creates an error for invalid data.
//
// Parameters: message — human-readable description.
// Returns: *InvalidDataError implementing error.
func NewInvalidDataError(message string) *InvalidDataError {
	return &InvalidDataError{
		abstractError{
			message: message,
		},
	}
}

// buildMessage returns message; if previous is non-nil, appends "; previous: ...".
func (e *abstractError) buildMessage(message string, previous error) string {
	err := message
	if previous != nil {
		err = fmt.Sprintf("%s; previous: %v", message, previous.Error())
	}
	return err
}

// SetPreviousError sets the wrapped cause and returns the receiver for chaining.
//
// Parameters: previous — the underlying error (may be nil).
// Returns: the receiver so callers can chain.
// Behavior: Error() will include "; previous: <previous.Error()>" when non-nil.
func (e *abstractError) SetPreviousError(previous error) *abstractError {
	e.previous = previous
	return e
}

// Error implements the error interface.
//
// Returns: the error message; if a previous error was set, it is appended.
func (e *abstractError) Error() string {
	return e.buildMessage(e.message, e.previous)
}
