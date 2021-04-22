package errorX

import "github.com/pkg/errors"

type Frame string

// StackTrace 將外部套件的型別, 改為我方系統的型別
type StackTrace []Frame

type pkgErrStack interface {
	StackTrace() errors.StackTrace
}

type pkgErrCause interface {
	Unwrap() error
	Error() string
}
