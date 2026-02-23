/*
Conjunto de erros customizados.

serve para padronizar todo lançamento de erro, pois estes são erros conhecidos
da aplicação.

É recomendável usar esta suíte de erros em toda a aplicação, pois sabemos como tratá-los.
*/
package domain

import "fmt"

type (
	abstractError struct {
		message  string
		previous error
	}
	NotFoundError struct {
		abstractError
	}
	InternalError struct {
		abstractError
	}
	ValidationError struct {
		abstractError
	}
	AlreadyExistsError struct {
		abstractError
	}
	DependencyError struct {
		abstractError
	}
	InvalidDataError struct {
		abstractError
	}
)

func NewDependencyError(message string) *DependencyError {
	return &DependencyError{
		abstractError{
			message: message,
		},
	}
}

func NewNotFoundError(message string) *NotFoundError {
	return &NotFoundError{
		abstractError{
			message: message,
		},
	}
}

func NewInternalError(message string) *InternalError {
	return &InternalError{
		abstractError{
			message: message,
		},
	}
}

func NewAlreadyExistsError(message string) *AlreadyExistsError {
	return &AlreadyExistsError{
		abstractError{
			message: message,
		},
	}
}

func NewInvalidDataError(message string) *InvalidDataError {
	return &InvalidDataError{
		abstractError{
			message: message,
		},
	}
}

// retorna a mensagem e caso haja um erro anterior, concatena na mensagem atual
func (e *abstractError) buildMessage(message string, previous error) string {
	err := message
	if previous != nil {
		err = fmt.Sprintf("%s; previous: %v", message, previous.Error())
	}
	return err
}

func (e *abstractError) SetPreviousError(previous error) *abstractError {
	e.previous = previous
	return e
}

func (e *abstractError) Error() string {
	return e.buildMessage(e.message, e.previous)
}
