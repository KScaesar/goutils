package errors

import (
	"net/http"
	"strings"

	"github.com/pkg/errors"
)

const UndefinedCode = -1

func Code(err error) int {
	if err == nil {
		return 0
	}

	if IsUndefinedError(err) {
		return UndefinedCode
	}

	return errors.Cause(err).(customError).Code()
}

func HttpStatus(err error) int {
	if err == nil {
		return http.StatusOK
	}

	if IsUndefinedError(err) {
		return http.StatusInternalServerError
	}

	return errors.Cause(err).(customError).HttpStatus()
}

// SimpleInfo 只取第一次 Wrap 的資訊, 後續 WrapMessage 函數的訊息不會保留,
// 如果 input 傳入的 err 就是 root error, 沒有經過 Wrap,
// 則只顯示 root error description
//
// Deprecated
func SimpleInfo(err error) string {
	var before, current, msgErr error

	current = err
	msgErr = current

	for current != nil {
		causeErrFactory, ok := current.(pkgErrCause)
		if !ok { // 沒找到就表示, current is root error
			msgList := strings.Split(msgErr.Error(), ": "+current.Error())
			return msgList[0]
		}

		before = current
		current = causeErrFactory.Unwrap()
		msgErr = before
	}

	return "no error, this all OK"
}

// Description 顯示初始錯誤 文字描述
func Description(err error) string {
	if err == nil {
		return "ok"
	}
	return errors.Cause(err).Error()
}
