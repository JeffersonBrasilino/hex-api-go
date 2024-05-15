package cqrs

type ActionHandler[TAction any] interface {
	Handle(data TAction) (any, error)
}
