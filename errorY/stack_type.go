package errorY

import (
	"github.com/pkg/errors"
)

// Frame is data from errors.Frame
type Frame string

// StackTrace is data from pkgErrStack
type StackTrace []Frame

type pkgErrStack interface {
	StackTrace() errors.StackTrace
}

type pkgErrCause interface {
	Unwrap() error
	Error() string
}
