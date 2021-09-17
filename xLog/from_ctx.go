package xLog

import (
	"context"
)

// example:
// https://pkg.go.dev/github.com/rs/zerolog#Logger.WithContext
// https://pkg.go.dev/github.com/rs/zerolog#Logger.UpdateContext

type logKey struct{}

func ContextWithLogger(ctx context.Context, l WrapperLogger) (logCtx context.Context) {
	if ctx == nil {
		ctx = context.Background()
	}
	return context.WithValue(ctx, logKey{}, &l)
}

func LoggerFromContext(logCtx context.Context) WrapperLogger {
	if logCtx == nil {
		return Logger()
	}

	logger, ok := logCtx.Value(logKey{}).(*WrapperLogger)
	if ok {
		return *logger
	}
	return Logger()
}
