package logY

import (
	"github.com/rs/zerolog"

	"github.com/Min-Feng/goutils/errors"
)

// SetGlobalLevel level  目前使用的字串有: "debug", "info", "error", "panic"
func SetGlobalLevel(level string) error {
	if level == "" {
		return errors.Wrap(errors.ErrInvalidParams, "level string is empty")
	}

	Level, err := zerolog.ParseLevel(level)
	if err != nil {
		return errors.Wrap(errors.ErrSystem, "parse log level, level=%v", level)
	}

	zerolog.SetGlobalLevel(Level)
	return nil
}

func IsDebugLevel() bool {
	return zerolog.GlobalLevel() == zerolog.DebugLevel
}
