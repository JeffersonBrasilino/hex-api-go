package core

type Result struct {
	success bool
	errors  any
	value   any
}

func ResultSuccess(val any) *Result {
	return &Result{success: true, errors: nil, value: val}
}

func ResultFailure(val any) *Result {
	return &Result{success: false, errors: val, value: nil}
}

func Combine(values ...any) *Result {
	errors := make([]error, len(values))
	return &Result{success: false, errors: errors, value: nil}
}

func (a *Result) IsSuccess() bool {
	return a.success
}

func (a *Result) GetValue() any {
	return a.value
}

func (a *Result) GetError() any {
	return a.errors
}
