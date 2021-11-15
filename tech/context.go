package tech

import (
	"context"

	"github.com/Min-Feng/goutils/xLog"
)

func SetRequestID(ctx context.Context, traceID string) context.Context {
	return xLog.ContextWithRequestID(ctx, traceID)
}

func RequestID(ctx context.Context) string {
	return xLog.RequestIDFromContext(ctx)
}

func SetLogger(ctx context.Context, l xLog.WrapperLogger) context.Context {
	return xLog.ContextWithLogger(ctx, l)
}

func Logger(ctx context.Context) xLog.WrapperLogger {
	return xLog.LoggerFromContext(ctx)
}
