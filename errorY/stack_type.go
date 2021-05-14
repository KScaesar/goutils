package errorY

import (
	"github.com/pkg/errors"
)

// following context, from custom type

type Frame string

type StackTrace []Frame

// following context, from pkg errors

type pkgErrStack interface {
	StackTrace() errors.StackTrace
}

type pkgErrCause interface {
	Unwrap() error
	Error() string
}
