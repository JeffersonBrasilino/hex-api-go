package endpoint

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/hex-api-go/pkg/core/infrastructure/message_system/message"
)

type serviceActivator struct {
	targetService any
	methodName    string
	methodToCall  reflect.Value
}

func NewServiceActivator(targetService any, methodName string) *serviceActivator {

	methodToCall := reflect.ValueOf(targetService).MethodByName(methodName)

	if !methodToCall.IsValid() {
		panic(fmt.Sprintf("method %s is invalid", methodName))
	}

	if methodToCall.Type().NumIn() != 1 ||
		methodToCall.Type().In(0).String() != "*message.Message" ||
		methodToCall.Type().NumOut() != 2 ||
		!(methodToCall.Type().Out(0).String() == "*message.Message" &&
			methodToCall.Type().Out(1).String() == "error") {

		panic(fmt.Sprintf("method %s is invalid, it must satisfy the signature of the MessageHandler.Handle method", methodName))

	}

	return &serviceActivator{
		targetService: targetService,
		methodName:    methodName,
		methodToCall:  methodToCall,
	}
}

func (s *serviceActivator) Handle(msg *message.Message) (*message.Message, error) {

	args := []reflect.Value{reflect.ValueOf(msg)}
	result := s.methodToCall.Call(args)

	resultMessage, ok := result[0].Interface().(*message.Message)
	if !ok {
		resultMessage = s.buildMessageFromPreviousMessage(msg, resultMessage)
	}

	err, okerr := result[1].Interface().(error)
	if !okerr {
		err = nil
	}

	if msg.GetHeaders().ReplyChannel != nil {
		msg.GetHeaders().ReplyChannel.Send(resultMessage)
	}
	
	return resultMessage, err
}

func (s *serviceActivator) buildMessageFromPreviousMessage(
	previousMessage *message.Message,
	payload any,
) *message.Message {
	payloadMarshaled, err := json.Marshal(payload)
	if err != nil {
		panic(fmt.Sprintf("can't marshal payload: %v", err))
	}
	return message.NewMessageBuilderFromMessage(previousMessage).
		WithPayload(payloadMarshaled).
		Build()
}
