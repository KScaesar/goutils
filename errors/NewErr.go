package errors

import (
	"github.com/cockroachdb/errors/errutil"
)

// New 只用在 global Err
func New(msg string) error {
	return errutil.NewWithDepth(1, msg)
}
