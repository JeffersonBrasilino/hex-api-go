package domain

import (
	"encoding/json"
	"regexp"

	"github.com/jeffersonbrasilino/ddgo"
)

var strongPasswordRegexes = []*regexp.Regexp{
	regexp.MustCompile(`[a-z]`),
	regexp.MustCompile(`[A-Z]`),
	regexp.MustCompile(`\d`),
	regexp.MustCompile(`[@$!%*?&]`),
	regexp.MustCompile(`^[A-Za-z\d@$!%*?&]{8,}$`),
}

type PasswordProps struct {
	Value string `domainValidator:"required"`
}

type Password struct {
	value string
}

func NewPassword(props *PasswordProps) (*Password, error) {
	if err := validatePassword(props); err != nil {
		return nil, err
	}

	return &Password{
		value: props.Value,
	}, nil
}

func validatePassword(props *PasswordProps) error {
	validator := ddgo.ValidatorInstance()
	validationErrors, faliedValidation := validator.Validate(props)
	if faliedValidation != nil {
		return ddgo.NewInternalError("Error when validating password data")
	}

	if len(validationErrors) > 0 {
		validationResult, failed := json.Marshal(validationErrors)
		if failed != nil {
			return ddgo.NewInternalError("Error when marshaling validation errors")
		}
		return ddgo.NewInvalidDataError(string(validationResult))
	}

	// Custom strong password validation using the required regex
	if props.Value != "" {
		isValid := true
		for _, re := range strongPasswordRegexes {
			if !re.MatchString(props.Value) {
				isValid = false
				break
			}
		}

		if !isValid {
			customErrors := map[string]any{
				"Value": map[string]any{
					"IsValid":          false,
					"FailedValidators": []string{"strong_password"},
				},
			}
			validationResult, _ := json.Marshal(customErrors)
			return ddgo.NewInvalidDataError(string(validationResult))
		}
	}

	return nil
}

func (p *Password) Value() string {
	return p.value
}
