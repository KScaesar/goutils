package logger

import (
	"github.com/rs/zerolog"
)

var (
	default_ WrapperLogger
)

// Logger 提供一個快速使用的函數, 固定輸出到 os.Stdout,
// 若希望同時輸出到多個 io.Writer, 此函數不適合, 請使用 New
func Logger() WrapperLogger {
	return default_
}

func Debug() *zerolog.Event {
	return default_.Debug().Caller(1)
}

func Info() *zerolog.Event {
	return default_.Info().Caller(1)
}

func Err(err error) *zerolog.Event {
	return default_.Err(err).Caller(1)
}

func Panic() *zerolog.Event {
	return default_.Panic().Caller(1)
}
