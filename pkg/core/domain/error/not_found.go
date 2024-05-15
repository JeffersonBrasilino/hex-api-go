package error

type NotFound struct {
	Message string
}

func (e *NotFound) Error() string {
	return e.Message
}
