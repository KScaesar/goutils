package errors

import (
	"github.com/pkg/errors"
)

// Frame is data from errors.Frame
type Frame string

// Stack is data from errors.StackTrace
type Stack []Frame

type pkgErrStack interface {
	StackTrace() errors.StackTrace
}

type pkgErrCause interface {
	Unwrap() error
	Error() string
}
