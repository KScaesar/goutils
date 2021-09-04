package errors

import (
	"github.com/pkg/errors"
)

// NewPkgError 產生 error 的時候, 也保存了呼叫地點的 stack 訊息, 若在全域變數呼叫, stack 訊息毫無價值
func NewPkgError(msg string) error {
	return errors.New(msg)
}

// Wrap 等效 WrapMessage + WrapStack,
// 最好在底層元件使用, 然後往上傳遞,
// error 傳遞的過程, 儘量不再使用 Wrap
func Wrap(err error, msg string, args ...interface{}) error {
	return errors.Wrapf(err, msg, args...)
}

// WrapMessage 通常用在 error 傳遞過程途中,
// 增加額外文字描述, 不會附加 stack 訊息,
// 第一個參數 err 應該在底層已經通過 Wrap 的轉換
func WrapMessage(err error, msg string, args ...interface{}) error {
	return errors.WithMessagef(err, msg, args...)
}

// WrapStack
// 每次呼叫都會增加新的 error stack 訊息,
// stack 有多個 frame, 從上到下的呼叫鏈
func WrapStack(err error) error {
	return errors.WithStack(err)
}

func Is(err, target error) bool {
	return errors.Is(err, target)
}

func Cause(err error) error {
	return errors.Cause(err)
}

func As(err error, target interface{}) bool {
	return errors.As(err, target)
}
