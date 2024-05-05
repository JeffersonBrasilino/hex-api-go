package core

type baseError struct {
	message string
}

type NotFound struct {
	baseError
}

type Internal struct {
	baseError
}

type Validation struct {
	baseError
}

type AlreadyExists struct {
	baseError
}

func (e *NotFound) Error() string {
	return e.baseError.message
}

func (e *Internal) Error() string {
	return e.baseError.message
}

func (e *Validation) Error() string {
	return e.baseError.message
}

func (e *AlreadyExists) Error() string {
	return e.baseError.message
}