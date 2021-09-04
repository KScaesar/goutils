package errors

import "github.com/pkg/errors"

// New 不可隨意新增 自定義錯誤, 儘量搭配 Wrap 產生 error,
// 主要會在全域使用, 以方便進行 error assert Is(),
func New(code int, httpStatus int, description string) error {
	return &baseError{code: code, httpStatus: httpStatus, description: description}
}

type baseError struct {
	code        int
	httpStatus  int
	description string
}

func (b *baseError) Error() string {
	return b.description
}

func IsUndefinedError(err error) bool {
	_, ok := errors.Cause(err).(*baseError)
	return !ok
}
