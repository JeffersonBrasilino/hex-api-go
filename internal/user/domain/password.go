// Package domain provides the core domain entities and value objects for the user bounded context.
// This file defines the Password value object, responsible for encapsulating and enforcing
// strong password rules for user authentication.
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

// PasswordProps holds the input data required to create a Password value object.
// The Value field is mandatory and must satisfy the strong password policy.
type PasswordProps struct {
	Value string `domainValidator:"required"`
}

// Password is an immutable value object that represents a validated strong password.
type Password struct {
	value string
}

// NewPassword creates and returns a new Password value object from the provided PasswordProps.
// It validates the required field using the domain validator and then enforces the strong
// password policy: at least 8 characters, uppercase, lowercase, digit, and special character.
//
// Parameters:
//   - props: pointer to PasswordProps containing the raw password string to validate.
//
// Returns a pointer to a valid Password and nil error on success, or nil and a domain error
// if validation fails.
func NewPassword(props *PasswordProps) (*Password, error) {
	if err := validatePassword(props); err != nil {
		return nil, err
	}
	return &Password{value: props.Value}, nil
}

// NewPasswordFromHash creates a Password value object directly from a pre-computed hash string,
// bypassing strength-rule validation. Use exclusively when reconstituting a User from storage.
func NewPasswordFromHash(hash string) *Password {
	return &Password{value: hash}
}

// validatePassword runs structural and strong-password validations against the given PasswordProps.
//
// Parameters:
//   - props: pointer to PasswordProps to validate.
//
// Returns nil on success, or a domain error describing the validation failure.
func validatePassword(props *PasswordProps) error {
	validator := ddgo.ValidatorInstance()
	validationErrors, err := validator.Validate(props)
	if err != nil {
		return ddgo.NewInternalError("Error when validating password data")
	}

	if len(validationErrors) > 0 {
		validationResult, err := json.Marshal(validationErrors)
		if err != nil {
			return ddgo.NewInternalError("Error when marshaling validation errors")
		}
		return ddgo.NewInvalidDataError(string(validationResult))
	}

	if props.Value != "" {
		for _, re := range strongPasswordRegexes {
			if !re.MatchString(props.Value) {
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
	}

	return nil
}


// Value returns the validated password string held by the Password value object.
func (p *Password) Value() string {
	return p.value
}
