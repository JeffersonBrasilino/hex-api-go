package domain

type customHeaders (map[string]string)

type PayloadType interface {
	~string | ~[]byte | struct{}
}

type Action[TPayload any] interface {
	Payload() TPayload
	Headers() customHeaders
	Name() string
}

type Command[T any] interface {
	Action[T]
}
type Query[T any] interface {
	Action[T]
}
type Event[T any] interface {
	Action[T]
}
