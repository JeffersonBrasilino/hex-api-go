package domain

type NotFoundError struct {
	Message       string
}

type InternalError struct {
	Message  string
}

type ValidationError struct {
	Message  string
}

type AlreadyExistsError struct {
	Message  string
}

type DependencyError struct {
	Message  string
}

func (e *NotFoundError) Error() string {
	return e.Message
}

func (e *InternalError) Error() string {
	return e.Message
}

func (e *ValidationError) Error() string {
	return e.Message
}

func (e *AlreadyExistsError) Error() string {
	return e.Message
}

func (e *DependencyError) Error() string {
	return e.Message
}
