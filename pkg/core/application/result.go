package application

type Result[TSuccess any, TError any] struct {
	success bool
	errors  TError
	value   TSuccess
}

func NewResult[TSuccess any, TError any]() *Result[TSuccess, TError] {
	return &Result[TSuccess, TError]{}
}

func (r *Result[TSuccess, TError]) Success(value TSuccess) {
	r.success = true
	r.value = value
}

func (r *Result[TSuccess, TError]) Error(value TError) {
	r.success = false
	r.errors = value
}

func (a *Result[TSuccess, TError]) IsSuccess() bool {
	return a.success
}

func (a *Result[TSuccess, TError]) GetValue() TSuccess {
	return a.value
}

func (a *Result[TSuccess, TError]) GetError() TError {
	return a.errors
}
