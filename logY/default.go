package logY

import (
	"github.com/rs/zerolog"
)

var (
	_default WrapperLogger
)

// Logger 若希望多個用戶同時寫入,
// 不同的 io.Writer, 此函數不適合
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
