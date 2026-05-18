package apperror

import "fmt"

type Error struct {
	Status  int
	Message string
	Fields  map[string]string
	Err     error
}

func (e *Error) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

func New(status int, message string) *Error {
	return &Error{
		Status:  status,
		Message: message,
	}
}

func WithFields(status int, message string, fields map[string]string) *Error {
	return &Error{
		Status:  status,
		Message: message,
		Fields:  fields,
	}
}

func Wrap(status int, message string, err error) *Error {
	return &Error{
		Status:  status,
		Message: message,
		Err:     err,
	}
}
