package validator

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

const (
	tagName = "domainValidator"
)

var (
	regexValidatorParams = regexp.MustCompile("=.+")
	regexValidatorName   = regexp.MustCompile(".+=")
	validators           = map[string]validateFunc{
		"required": requiredValidator,
		"gte":      gteValidator,
		"lte":      lteValidator,
		"len":      lenValidator,
	}
)

type (
	fieldError      []any
	errorsResult    map[string]any
	validateFunc    func(reflect.Value, string) string
	domainValidator struct {
		valid      bool
		errors     errorsResult
		validators map[string]validateFunc
	}
)

/*
validador de dados de criação do dominio
*/
func NewDomainValidator() *domainValidator {
	return &domainValidator{
		valid:      true,
		errors:     make(errorsResult),
		validators: validators,
	}
}

// valida os dados de criação do dominio
func (d *domainValidator) Validate(data any) bool {
	values := reflect.ValueOf(data)
	if values.Kind() == reflect.Ptr {
		values = reflect.Indirect(values)
	}
	d.errors = d.validateProps(values)
	return d.valid
}

// retorna os erros de dados de criação do dominio
func (d *domainValidator) GetErrors() errorsResult {
	return d.errors
}

// adiciona uma função de validação customizada para o dado de criação do dominio
func (d *domainValidator) AddCustomValidator(name string, validator validateFunc) error {
	if d.validators[name] != nil {
		return fmt.Errorf("validator for name %s already exists", name)
	}
	d.validators[name] = validator
	return nil
}

func (d *domainValidator) validateProps(prop reflect.Value) errorsResult {
	errors := make(errorsResult)
	var dataSize int
	if prop.Kind() == reflect.Struct {
		dataSize = prop.NumField()
	} else {
		dataSize = prop.Len()
	}

	var value reflect.Value
	var propName string
	var tag string
	for i := 0; i < dataSize; i++ {
		if prop.Kind() == reflect.Struct {
			value = prop.Field(i)
			propName = prop.Type().Field(i).Name
			tag = prop.Type().Field(i).Tag.Get(tagName)
		} else {
			value = prop.Index(i)
			propName = strconv.Itoa(i)
		}

		if value.Kind() == reflect.Struct {
			errors[propName] = d.validateProps(value)
		}

		if tag == "" || tag == "-" {
			continue
		}
		res := d.executeValidatorsByTag(tag, value)
		if len(res) > 0 {
			errors[propName] = res
		}
	}
	return errors
}

func (d *domainValidator) executeValidatorsByTag(tag string, value reflect.Value) fieldError {
	errors := fieldError{}
	for _, tg := range strings.Split(tag, ",") {
		if tg == "nested" {
			res := d.validateProps(value)
			if len(res) > 0 {
				errors = append(errors, res)
			}
			continue
		}

		tagName := regexValidatorParams.ReplaceAllLiteralString(tg, "")
		if _, ok := validators[tagName]; !ok {
			panic(fmt.Sprintf("validator %s does not exist", tagName))
		}
		var tagParams string
		if regexValidatorParams.MatchString(tg) {
			tagParams = regexValidatorName.ReplaceAllLiteralString(tg, "")
		}
		err := validators[tagName](value, tagParams)
		if err != "" {
			d.valid = false
			errors = append(errors, err)
		}
	}
	return errors
}
