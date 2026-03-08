// Validators used by the domain validator (tag-based validation).
//
// Intent: Provide required, gte, lte, and len checks for struct fields
// when using the domainValidator tag.
// Objective: Keep validation logic in one place and reusable by
// domain_validator without exporting these helpers.
package domain

import (
	"reflect"
	"strconv"
)

func isComplexValue(value reflect.Value) bool {
	return value.Kind() == reflect.Interface ||
		value.Kind() == reflect.Struct ||
		value.Kind() == reflect.Slice ||
		value.Kind() == reflect.Map
}

func requiredValidator(params string) validateFunc {
	return func(value reflect.Value) bool {
		if isComplexValue(value) && value.IsNil() {
			return false
		}
		if value.IsZero() {
			return false
		}
		return true
	}
}

func gteValidator(param string) validateFunc {
	return func(value reflect.Value) bool {
		len, err := strconv.Atoi(param)
		if err != nil {
			return false
		}
		if value.Len() < len {
			return false
		}
		return true
	}
}

func lteValidator(param string) validateFunc {
	return func(value reflect.Value) bool {
		len, err := strconv.Atoi(param)
		if err != nil {
			return false
		}
		if value.Len() > len {
			return false
		}
		return true
	}
}

func lenValidator(param string) validateFunc {
	return func(value reflect.Value) bool {
		len, err := strconv.Atoi(param)
		if err != nil {
			return false
		}
		if value.Len() == len {
			return false
		}
		return true
	}
}
