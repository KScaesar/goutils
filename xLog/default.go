package xLog

import (
	"github.com/rs/zerolog"
)

var (
	default_ WrapperLogger
)

func SetDefaultLogger(l WrapperLogger) {
	default_ = l
}

// Logger 提供一個快速使用的函數, 預設輸出到 os.Stdout,
// 若希望同時輸出到多個 io.Writer, 此函數不適合
func Logger() WrapperLogger {
	return default_
}

func Debug() *zerolog.Event {
	return default_.Unwrap().Debug().Caller(1)
}

func Info() *zerolog.Event {
	return default_.Unwrap().Info().Caller(1)
}

func Err(err error) *zerolog.Event {
	return default_.Unwrap().Err(err).Caller(1)
}

func Error() *zerolog.Event {
	return default_.Unwrap().Error().Caller(1)
}

func Panic() *zerolog.Event {
	return default_.Unwrap().Panic().Caller(1)
}
