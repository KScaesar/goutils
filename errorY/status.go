package errorY

import (
	"net/http"

	"github.com/pkg/errors"
)

func Code(err error) int {
	base, ok := errors.Cause(err).(*baseError)
	if !ok {
		return -1
	}
	return base.code
}

func HTTPStatus(err error) int {
	base, ok := errors.Cause(err).(*baseError)
	if !ok {
		return http.StatusInternalServerError
	}
	return base.httpStatus
}

// SimpleInfo 只取第一次 Wrap 的資訊, 後續 WrapMessage 函數的訊息不會保留,
// 如果一開始傳入的 err 就是 root error,
// 則只顯示 root error descrption
func SimpleInfo(err error) string {
	var before, current, msgErr error

	current = err
	msgErr = current

	for current != nil {
		causeErrFactory, ok := current.(pkgErrCause)
		if !ok { // 沒找到就表示, current is root error
			return msgErr.Error()
		}

		before = current
		current = causeErrFactory.Unwrap()
		msgErr = before
	}
	return "no error, this all OK"
}

// Description 顯示初始錯誤 文字描述
func Description(err error) string {
	return errors.Cause(err).Error()
}
