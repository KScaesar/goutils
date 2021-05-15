package errorY

import (
	"github.com/pkg/errors"
)

// the following, from custom type

type Frame string

type StackTrace []Frame

// the following, from pkg errors

type pkgErrStack interface {
	StackTrace() errors.StackTrace
}

type pkgErrCause interface {
	Unwrap() error
	Error() string
}
