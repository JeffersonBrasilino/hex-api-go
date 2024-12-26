package endpoint

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"

	"github.com/hex-api-go/pkg/core/infrastructure/message_system/message"
)

type serviceActivator struct {
	methodToCall reflect.Value
	paramInputs  []reflect.Type
	paramOutputs []reflect.Type
}

func NewServiceActivator(targetService any, methodName string) *serviceActivator {

	methodToCall := reflect.ValueOf(targetService).MethodByName(methodName)
	inputs := []reflect.Type{}
	for i := 0; i < methodToCall.Type().NumIn(); i++ {
		var p reflect.Type
		if methodToCall.Type().In(i).Kind() != reflect.Ptr {
			p = methodToCall.Type().In(i)
		} else {
			p = methodToCall.Type().In(i).Elem()
		}
		inputs = append(inputs, p)
	}

	outputs := []reflect.Type{}
	for i := 0; i < methodToCall.Type().NumOut(); i++ {
		var p reflect.Type
		if methodToCall.Type().Out(i).Kind() != reflect.Ptr {
			p = methodToCall.Type().Out(i)
		} else {
			p = methodToCall.Type().Out(i).Elem()
		}
		outputs = append(outputs, p)
	}

	fmt.Println("OUTPUT PARAM", methodToCall.Type().NumOut())

	return &serviceActivator{
		methodToCall: methodToCall,
		paramInputs:  inputs,
		paramOutputs: outputs,
	}
}

func (s *serviceActivator) Handle(msg *message.Message) (*message.Message, error) {

	var args []reflect.Value
	var err error
	sizeInputParam := len(s.paramInputs)
	if sizeInputParam >= 1 {
		args, err = s.makeInputs(msg.GetPayload())
		if err != nil {
			return nil, err
		}
	}

	resultCall := s.methodToCall.Call(args)
	var result []any
	for _, v := range resultCall {
		result = append(result, v.Interface())
	}

	resultMarshal, errMsrl := json.Marshal(result)
	if errMsrl != nil {
		return nil, fmt.Errorf("[service-activator] cannot marshal response: %s", errMsrl)
	}

	resultMessage := message.NewMessageBuilder().
		WithChannelName(msg.GetHeaders().ReplyChannel.Name()).
		WithMessageType(message.Document).
		WithPayload(resultMarshal).
		Build()

	resultMessage.SetInternalPayload(result)
	if msg.GetHeaders().ReplyChannel != nil {
		msg.GetHeaders().ReplyChannel.Send(resultMessage)
	}

	return resultMessage, err
}

func (s *serviceActivator) makeInputs(data []byte) ([]reflect.Value, error) {

	var err error

	var dataParsed = map[string]any{}
	if len(s.paramInputs) > 1 {
		err = json.Unmarshal(data, &dataParsed)
		if err != nil {
			err = fmt.Errorf("[service-activator] %s", err)
			return nil, err
		}
	} else {

		var allDataMap = map[string]any{}
		err = json.Unmarshal(data, &allDataMap)
		if err != nil {
			err = fmt.Errorf("[service-activator] %s", err)
			return nil, err
		}
		dataParsed["0"] = allDataMap

	}

	var args = []reflect.Value{}
	dataKey := 0
	for _, v := range dataParsed {
		if s.isComplexValue(s.paramInputs[dataKey]) {
			value, errCv := s.makeComplexValue(s.paramInputs[dataKey], v)
			if errCv != nil {
				err = errCv
				break
			}
			args = append(args, value)
			continue
		}

		value, errCv := s.makePrimitiveValue(s.paramInputs[dataKey], v)
		if errCv != nil {
			err = errCv
			break
		}
		args = append(args, value)
		dataKey++
	}
	return args, err
}

func (s *serviceActivator) makeComplexValue(field reflect.Type, value any) (reflect.Value, error) {
	var err error
	if field.Kind() == reflect.Struct || field.Kind() == reflect.Map {
		fieldResult := reflect.New(field).Interface()
		byteData, isByteArr := value.([]byte)
		if !isByteArr {
			byteData, err = json.Marshal(value)
		}

		if err != nil {
			return reflect.ValueOf(nil), fmt.Errorf("[service-activator] %s", err)
		}

		json.Unmarshal(byteData, fieldResult)
		return reflect.ValueOf(fieldResult), nil
	}

	var fieldResult = reflect.New(field).Elem()
	fieldResult.Elem()
	fieldResult.Set(reflect.ValueOf(value))
	fieldResult.Convert(field)
	return fieldResult, nil

}

func (s *serviceActivator) isComplexValue(value reflect.Type) bool {
	return value.Kind() == reflect.Interface ||
		value.Kind() == reflect.Struct ||
		value.Kind() == reflect.Slice ||
		value.Kind() == reflect.Map
}

func (s *serviceActivator) makePrimitiveValue(field reflect.Type, value any) (reflect.Value, error) {
	var fieldResult = reflect.New(field).Elem()
	var parsedValue, _ = value.(string)
	switch field.Kind() {
	case reflect.Bool:
		dt, _ := strconv.ParseBool(parsedValue)
		fieldResult.SetBool(dt)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int64:
		dt, _ := strconv.Atoi(parsedValue)
		fieldResult.SetInt(int64(dt))
	case reflect.Uint, reflect.Uint8, reflect.Uint64:
		dt, _ := strconv.Atoi(parsedValue)
		fieldResult.SetUint(uint64(dt))
	case reflect.Float32, reflect.Float64:
		dt, _ := strconv.ParseFloat(parsedValue, 64)
		fieldResult.SetFloat(dt)
	case reflect.String:
		fieldResult.SetString(parsedValue)
	default:
		return reflect.ValueOf(nil), fmt.Errorf("[service-activator] invalid value type")
	}

	return fieldResult, nil
}
