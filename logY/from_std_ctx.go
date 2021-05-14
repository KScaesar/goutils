package logY

import (
	"context"
)

type logKey struct{}

// example:
// https://pkg.go.dev/github.com/rs/zerolog#Logger.WithContext
// https://pkg.go.dev/github.com/rs/zerolog#Logger.UpdateContext

func FromCtx(ctx context.Context) WrapperLogger {
	logger, ok := ctx.Value(logKey{}).(*WrapperLogger)
	if ok {
		return *logger
	}
	return Logger()
}

// WithCtx 不要誤用 zerolog.Logger.WithContext, 通常會使用 自定義的 WrapperLogger.WithCtx
func (l WrapperLogger) WithCtx(ctx context.Context) context.Context {
	return context.WithValue(ctx, logKey{}, &l)
}
