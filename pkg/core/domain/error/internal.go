package error

type Internal struct {
	Message string
}

func (e *Internal) Error() string {
	return e.Message
}
