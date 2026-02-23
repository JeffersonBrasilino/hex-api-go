package domain

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"sync"
)

const (
	mainTagName = "domainValidator"
)

var (
	regexValidatorParams = regexp.MustCompile("=.+")
	regexValidatorName   = regexp.MustCompile(".+=")
	validatorInstance    = &domainValidator{
		cache: make(map[reflect.Type][]validationStep),
	}
	validators = map[string]func(params string) validateFunc{
		"required": requiredValidator,
		"gte":      gteValidator,
		"lte":      lteValidator,
		"len":      lenValidator,
	}
)

type validateFunc func(value reflect.Value) bool

type validateResult struct {
	IsValid          bool
	FailedValidators []string
}

type domainValidator struct {
	mu    sync.RWMutex
	cache map[reflect.Type][]validationStep
}

type validationStep struct {
	fieldName        string
	indexField       []int
	validators       map[string]validateFunc
	valueExtractFunc func(reflect.Value) reflect.Value
}

type parentFieldData struct {
	parentIndexPosition []int
	parentFieldName     []string
}

// ValidatorInstance returns the singleton domain validator used across the
// package. It never returns nil.
//
// Intent: provide a shared validator instance that caches validation schema
// definitions by type to avoid repeated reflection work.
//
// Returns: pointer to `domainValidator` singleton.
func ValidatorInstance() *domainValidator {
	return validatorInstance
}

// ResetDomainValidatorForTest resets the package validator singleton. This
// function is intended to be used in tests to clear cached schemas and return
// the validator to a clean state.
func ResetDomainValidatorForTest() {
	validatorInstance = &domainValidator{
		cache: make(map[reflect.Type][]validationStep),
	}
}

// Validate inspects the provided struct value and runs the configured domain
// validators declared on struct fields using the `domainValidator` tag.
//
// Intent: perform validation of a DTO or domain input using reflection and a
// cached validation schema to reduce repeated reflection cost.
//
// Parameters:
//   - data: any struct pointer or value to validate. If nil, the function
//     returns (nil, nil).
//
// Returns:
//   - map[string]validateResult: a map keyed by the dotted field path
//     (e.g. "Address.Street") containing validation outcomes for failing
//     fields. If the map is empty, no validation errors were found.
//   - error: non-nil if an internal error occurred while building the schema.
//
// Behavior:
//   - The function builds or retrieves a cached schema for the concrete type
//     and executes each configured validator. It does not short-circuit; all
//     configured validators for each field are evaluated and any failures are
//     reported in the returned map.
//
// Example usage:
//
//	type createUserDTO struct {
//	    Name string `domainValidator:"required"`
//	    Age  int    `domainValidator:"gte=18"`
//	}
//
//	dto := createUserDTO{Name: "", Age: 17}
//	result, _ := ValidatorInstance().Validate(&dto)
//	// result will contain validation failures for Name and Age.
//	available validators: required, gte, lte, len
func (d *domainValidator) Validate(data any) (map[string]validateResult, error) {
	if data == nil {
		return nil, nil
	}

	values := reflect.Indirect(reflect.ValueOf(data))
	validationSchema, err := d.getValidationSchema(data)

	if err != nil {
		return nil, err
	}

	result := make(map[string]validateResult)
	for _, step := range validationSchema {
		value := step.valueExtractFunc(values)
		faliedValidatorsTag := []string{}
		for tag, validatorFunc := range step.validators {
			isValid := validatorFunc(value)
			if isValid != false {
				continue
			}
			faliedValidatorsTag = append(faliedValidatorsTag, tag)
		}

		if len(faliedValidatorsTag) == 0 {
			continue
		}

		result[step.fieldName] = validateResult{
			IsValid:          false,
			FailedValidators: faliedValidatorsTag,
		}

	}

	return result, nil
}

func (d *domainValidator) getValidationSchema(data any) ([]validationStep, error) {
	key := reflect.TypeOf(data)
	if key.Kind() == reflect.Ptr {
		key = key.Elem()
	}

	d.mu.RLock()
	schema, ok := d.cache[key]
	d.mu.RUnlock()

	if !ok {
		d.mu.Lock()
		values := reflect.Indirect(reflect.ValueOf(data))
		newSchema, err := d.makeValidationSchema(values, nil)
		schema = newSchema

		if err != nil {
			d.mu.Unlock()
			return nil, err
		}
		d.cache[key] = schema
		d.mu.Unlock()
	}

	return schema, nil

}

func (d *domainValidator) makeValidationSchema(
	data reflect.Value,
	parentField *parentFieldData,
) ([]validationStep, error) {

	schemaValue := data
	fieldName := []string{}
	dataKeys := []int{}

	if parentField != nil {
		fieldName = parentField.parentFieldName
		schemaValue = data.FieldByIndex(parentField.parentIndexPosition)
		dataKeys = parentField.parentIndexPosition
	}

	if schemaValue.Kind() == reflect.Ptr {
		schemaValue = reflect.Indirect(schemaValue)
	}

	var dataSize int
	if schemaValue.Kind() == reflect.Struct {
		dataSize = schemaValue.NumField()
	} else {
		dataSize = schemaValue.Len()
	}

	result := []validationStep{}
	for i := 0; i < dataSize; i++ {
		currentPath := make([]int, len(dataKeys))
		copy(currentPath, dataKeys)
		currentPath = append(currentPath, i)
		fieldInfo := data.Type().FieldByIndex(currentPath)
		fieldName := append(fieldName, fieldInfo.Name)

		valueData := data.FieldByIndex(currentPath)
		if valueData.Kind() == reflect.Ptr {
			valueData = valueData.Elem()
		}

		if valueData.Kind() == reflect.Struct {

			childSteps, err := d.makeValidationSchema(data, &parentFieldData{
				parentIndexPosition: currentPath,
				parentFieldName:     fieldName,
			})

			if err != nil {
				return nil, err
			}

			result = append(result, childSteps...)
			continue
		}

		validatorsTags := fieldInfo.Tag.Get(mainTagName)
		if validatorsTags == "" {
			continue
		}
		validators, errValidators := d.makeValidationsFuncs(validatorsTags)
		if errValidators != nil {
			return nil, fmt.Errorf("[domain-validator] %s", errValidators.Error())
		}

		step := validationStep{
			fieldName:        strings.Join(fieldName, "."),
			indexField:       currentPath,
			validators:       validators,
			valueExtractFunc: d.makeStructDataAccessor(currentPath),
		}
		result = append(result, step)
	}
	return result, nil
}

func (d *domainValidator) makeValidationsFuncs(
	validationTags string,
) (map[string]validateFunc, error) {
	if validationTags == "" {
		return nil, fmt.Errorf(
			"validation tags is empty, add validation tags after tag domainValidator.",
		)
	}

	validatorsNames := strings.Split(validationTags, ",")
	validatorsArr := make(map[string]validateFunc, len(validatorsNames))
	for _, tag := range validatorsNames {
		tagName := regexValidatorParams.ReplaceAllLiteralString(tag, "")
		validator, ok := validators[tagName]
		if !ok {
			return nil, fmt.Errorf("validator %s does not exist", tagName)
		}

		var tagParams string
		if regexValidatorParams.MatchString(tag) {
			tagParams = regexValidatorName.ReplaceAllLiteralString(tag, "")
		}
		validatorsArr[tag] = validator(tagParams)
	}

	return validatorsArr, nil
}

func (d *domainValidator) makeStructDataAccessor(
	fieldAddress []int,
) func(dataStruct reflect.Value) reflect.Value {
	return func(dataStruct reflect.Value) reflect.Value {
		return dataStruct.FieldByIndex(fieldAddress)
	}
}