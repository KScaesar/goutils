package xLog

import (
	"context"
)

type logKey struct{}

func ContextWithLogger(ctx context.Context, l WrapperLogger) (logCtx context.Context) {
	return context.WithValue(ctx, logKey{}, &l)
}

func LoggerFromContext(logCtx context.Context) WrapperLogger {
	logger, ok := logCtx.Value(logKey{}).(*WrapperLogger)
	if !ok {
		l := Logger().Unwrap().With().Str("LoggerFromContext", "not exist").Logger()
		return WrapPrototype(&l)
	}
	return *logger
}
