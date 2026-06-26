// Package domain provides the core domain entities and value objects for the user bounded context.
// This file defines the Document value object, responsible for encapsulating and validating
// a Brazilian CPF document number.
package domain

import (
	"encoding/json"

	"github.com/jeffersonbrasilino/ddgo"
)

// DocumentProps holds the input data required to create a Document value object.
// The Value field is mandatory and must contain a non-empty CPF string.
type DocumentProps struct {
	Value string `domainValidator:"required"`
}

// Document is an immutable value object that represents a validated CPF document number.
type Document struct {
	value string
}

// NewDocument creates and returns a new Document value object from the provided DocumentProps.
// It validates the props using the domain validator and verifies the CPF checksum algorithm.
//
// Parameters:
//   - props: pointer to DocumentProps containing the raw CPF string to validate.
//
// Returns a pointer to a valid Document and nil error on success, or nil and a domain error
// if the required validation or CPF checksum fails.
func NewDocument(props *DocumentProps) (*Document, error) {
	if err := validateDocument(props); err != nil {
		return nil, err
	}
	return &Document{value: props.Value}, nil
}

// validateDocument runs structural and CPF checksum validations against the given DocumentProps.
//
// Parameters:
//   - props: pointer to DocumentProps to validate.
//
// Returns nil on success, or a domain error describing the validation failure.
func validateDocument(props *DocumentProps) error {
	validator := ddgo.ValidatorInstance()
	validationErrors, err := validator.Validate(props)
	if err != nil {
		return ddgo.NewInternalError("Error when validating contact data")
	}

	if len(validationErrors) > 0 {
		validationResult, err := json.Marshal(validationErrors)
		if err != nil {
			return ddgo.NewInternalError("Error when marshaling validation errors")
		}
		return ddgo.NewInvalidDataError(string(validationResult))
	}

	if !isValidCPF(props.Value) {
		return ddgo.NewInvalidDataError("Invalid cpf.")
	}

	return nil
}

// isValidCPF verifies whether the given string is a valid Brazilian CPF number.
// It checks the length, rejects sequences of identical digits, and validates both check digits.
//
// Parameters:
//   - value: the CPF string to validate (must contain exactly 11 numeric characters).
//
// Returns true if the CPF is valid, false otherwise.
func isValidCPF(value string) bool {
	if len(value) != 11 {
		return false
	}

	allSame := true
	for i := 1; i < 11; i++ {
		if value[i] != value[0] {
			allSame = false
			break
		}
	}
	if allSame {
		return false
	}

	cpf := make([]int, 11)
	for i, c := range value {
		cpf[i] = int(c - '0')
	}

	rest := func(count int) int {
		end := len(cpf) + count - 12
		sum := 0
		for index, el := range cpf[:end] {
			sum += el * (count - index)
		}
		return ((sum * 10) % 11) % 10
	}

	return rest(10) == cpf[9] && rest(11) == cpf[10]
}

// Value returns the validated CPF string held by the Document value object.
func (d *Document) Value() string {
	return d.value
}
