package errors

import "errors"

func NewStdError(msg string) error {
	return errors.New(msg)
}
