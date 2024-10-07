package validator

import (
	"fmt"
	"reflect"
	"strconv"
)

func isComplexValue(value reflect.Value) bool {
	return value.Kind() == reflect.Interface ||
		value.Kind() == reflect.Struct ||
		value.Kind() == reflect.Slice ||
		value.Kind() == reflect.Map
}

func requiredValidator(value reflect.Value, param string) string {
	if isComplexValue(value) {
		if value.IsNil() {
			return "is required"
		}
	}
	if value.IsZero() {
		return "is required"
	}
	return ""
}

func gteValidator(value reflect.Value, param string) string {
	len, err := strconv.Atoi(param)
	if err != nil {
		panic(fmt.Sprintf("invalid param value %v for gte validator", param))
	}
	if value.Len() < len {
		return fmt.Sprintf("must be more than %v characters", param)
	}
	return ""
}

func lteValidator(value reflect.Value, param string) string {
	len, err := strconv.Atoi(param)
	if err != nil {
		panic(fmt.Sprintf("invalid param value %v for lte validator", param))
	}
	if value.Len() > len {
		return fmt.Sprintf("must be less than %v characters", param)
	}
	return ""
}

func lenValidator(value reflect.Value, param string) string {
	len, err := strconv.Atoi(param)
	if err != nil {
		panic(fmt.Sprintf("invalid param value %v for lte validator", param))
	}
	if value.Len() > len {
		return fmt.Sprintf("must be less than %v characters", param)
	}
	return ""
}
