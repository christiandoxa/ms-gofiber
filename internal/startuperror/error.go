package startuperror

import (
	"errors"
	"fmt"
)

// Code identifies a startup failure phase.
type Code string

const (
	ConfigLoad    Code = "startup.config_load"
	AppBuild      Code = "startup.app_build"
	FiberListen   Code = "startup.fiber_listen"
	FiberShutdown Code = "startup.fiber_shutdown"
	AppClose      Code = "startup.app_close"
)

// Error wraps a startup failure with a stable code for logs and docs.
type Error struct {
	Code Code
	Err  error
}

func Wrap(code Code, err error) error {
	if err == nil {
		return nil
	}
	return &Error{Code: code, Err: err}
}

func (e *Error) Error() string {
	return fmt.Sprintf("%s: %v", e.Code, e.Err)
}

func (e *Error) Unwrap() error {
	return e.Err
}

func CodeOf(err error) (Code, bool) {
	if startupErr, ok := errors.AsType[*Error](err); ok {
		return startupErr.Code, true
	}
	return "", false
}
