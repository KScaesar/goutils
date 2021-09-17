package tech

import (
	"context"

	"github.com/Min-Feng/goutils/xLog"
)

func SetTraceID(ctx context.Context, traceID string) context.Context {
	return xLog.ContextWithTraceID(ctx, traceID)
}

func TraceID(ctx context.Context) string {
	return xLog.TraceIDFromContext(ctx)
}

func SetLogger(ctx context.Context, l xLog.WrapperLogger) context.Context {
	return xLog.ContextWithLogger(ctx, l)
}

func Logger(ctx context.Context) xLog.WrapperLogger {
	return xLog.LoggerFromContext(ctx)
}
