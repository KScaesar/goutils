package errorX

import "github.com/pkg/errors"

// Wrapf
// 第一個參數 err 應該傳入 global defined error, example: ErrSystemFailed
//
// Wrap 只能在底層元件使用, 然後往上傳遞,
// error 傳遞的過程, 儘量不再使用 Wrap 系列的函數
//
// 實際用法可以參考 TestStacks_CountOK
func Wrapf(err error, format string, args ...interface{}) error {
	return errors.Wrapf(err, format, args...)
}

func Wrap(err error, msg string) error {
	return errors.Wrap(err, msg)
}

// WithMsgf
// 第一個參數 err 應該在底層已經通過 Wrap 的轉換
// 通常用在 error 傳遞時, 增加額外描述
// 外部套件的錯誤, 不應該使用 WithMsg 系列的函數重新包裝
func WithMsgf(err error, msg string, args ...interface{}) error {
	return errors.WithMessagef(err, msg, args...)
}

func WithMsg(err error, msg string) error {
	return errors.WithMessage(err, msg)
}

func Is(err, target error) bool {
	return errors.Is(err, target)
}

// Cause 獲得自定義(baseError)的頂級錯誤, 而不是包裝(Wrap)過的 error
func Cause(err error) error {
	return errors.Cause(err)
}

func As(err error, target interface{}) bool {
	return errors.As(err, target)
}
