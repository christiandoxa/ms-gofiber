package apperror

import "fmt"

type Code string

const (
	ErrBadRequest   Code = "BAD_REQUEST"
	ErrValidation   Code = "VALIDATION"
	ErrUnauthorized Code = "UNAUTHORIZED"
	ErrForbidden    Code = "FORBIDDEN"
	ErrNotFound     Code = "NOT_FOUND"
	ErrConflict     Code = "CONFLICT"
	ErrDB           Code = "DB_ERROR"
	ErrInternal     Code = "INTERNAL"
)

type Error struct {
	Code    Code
	Message string
	Fields  map[string]string
	Err     error
}

func (e *Error) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func New(code Code, msg string) *Error {
	return &Error{Code: code, Message: msg}
}

func WithFields(code Code, msg string, fields map[string]string) *Error {
	return &Error{Code: code, Message: msg, Fields: fields}
}

func Wrap(code Code, msg string, err error) *Error {
	return &Error{Code: code, Message: msg, Err: err}
}
