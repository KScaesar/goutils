package logY

import (
	"github.com/Min-Feng/goutils/errorY"

	"github.com/rs/zerolog"
)

// SetGlobalLevel level  目前使用的字串有: "debug", "info", "error", "panic"
func SetGlobalLevel(level string) error {
	if level == "" {
		return errorY.Wrap(errorY.ErrInvalidParams, "level string is empty")
	}

	Level, err := zerolog.ParseLevel(level)
	if err != nil {
		return errorY.Wrap(errorY.ErrSystem, "parse log level, level=%v", level)
	}

	zerolog.SetGlobalLevel(Level)
	return nil
}

func IsDebugLevel() bool {
	return zerolog.GlobalLevel() == zerolog.DebugLevel
}
