/*
httpValidateRequest componente que auxilia a validação de dados de requests HTTP.

Usa como base o pacote [github.com/go-playground/validator/v10].

Todas as validações disponíveis olhar o pacote [github.com/go-playground/validator/v10].
*/
package http

import (
	"fmt"
	"regexp"

	"github.com/go-playground/validator/v10"
)

var (
	validate          *validator.Validate
	regexRequestField = regexp.MustCompile("[Rr]equest.")
)

type validatorResponse []string

/*
valida a estrutura(struct) da request.

ex:

	type Request struct {
		Username string `validate:"gte=4"`
		Password string `validate:"required"`
		Devices  Device `validate:"required"`
	}

	ValidateRequest(Request{})
*/
func ValidateRequest(requestData interface{}) validatorResponse {
	validate = validator.New()
	errs := validate.Struct(requestData)
	return parseErrors(errs)
}

/*
adiciona um validador customizado para a propriedade.

ex:

	type Request struct {
		Username string `validate:"customValidator"`
	}

	AddCustomValidator("customValidator", func(fl validator.FieldLevel) bool){
		//...
	})
*/
func AddCustomValidator(validatorName string, validatorFunc func(fl validator.FieldLevel) bool) {
	err := validate.RegisterValidation(validatorName, validatorFunc)
	if err != nil {
		panic(err)
	}
}

// normaliza os erros de validações
func parseErrors(err error) validatorResponse {

	errors := validatorResponse{}

	for _, err := range err.(validator.ValidationErrors) {
		fieldName := regexRequestField.ReplaceAllString(err.Namespace(), "")
		if err.Param() != "" {
			errors = append(errors, fmt.Sprintf("%s is %s to %s", fieldName, err.Tag(), err.Param()))
			continue
		}
		errors = append(errors, fmt.Sprintf("%s is %s", fieldName, err.Tag()))
	}

	return errors
}
