package logY

import (
	"github.com/rs/zerolog"
)

var (
	_default WrapperLogger
)

// Logger 提供一個快速使用的函數, 固定輸出到 os.Stdout,
// 若希望同時輸出到多個 io.Writer, 此函數不適合, 請使用 New
func Logger() WrapperLogger {
	return _default
}

func Debug() *zerolog.Event {
	return _default.Prototype().Debug().Caller(1)
}

func Info() *zerolog.Event {
	return _default.Prototype().Info().Caller(1)
}

func Error(err error) *zerolog.Event {
	return _default.ErrCode(err).Prototype().Err(err).Caller(1)
}

func Panic() *zerolog.Event {
	return _default.Prototype().Panic().Caller(1)
}
