/*
Conjunto de erros customizados.

serve para padronizar todo lançamento de erro, pois estes são erros conhecidos
da aplicação.

É recomendável usar esta suíte de erros em toda a aplicação, pois sabemos como tratá-los.
*/
package domain

import "fmt"

type (
	NotFoundError struct {
		Message  string
		Previous error
	}
	InternalError struct {
		Message  string
		Previous error
	}
	ValidationError struct {
		Message  string
		Previous error
	}
	AlreadyExistsError struct {
		Message  string
		Previous error
	}
	DependencyError struct {
		Message  string
		Previous error
	}
)

// retorna a mensagem e caso haja um erro anterior, concatena na mensagem atual
func buildMessage(message string, previous error) string {
	err := message
	if previous != nil {
		err = fmt.Sprintf("%s; previous: %v", message, previous.Error())
	}
	return err
}

func (e *NotFoundError) Error() string {
	return buildMessage(e.Message, e.Previous)
}

func (e *InternalError) Error() string {
	return buildMessage(e.Message, e.Previous)
}

func (e *ValidationError) Error() string {
	return buildMessage(e.Message, e.Previous)
}

func (e *AlreadyExistsError) Error() string {
	return buildMessage(e.Message, e.Previous)
}

func (e *DependencyError) Error() string {
	return buildMessage(e.Message, e.Previous)
}
