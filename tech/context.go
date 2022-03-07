package tech

import (
	"context"

	"github.com/KScaesar/goutils"
	"github.com/KScaesar/goutils/xLog"
)

func SetCorrelationID(ctx context.Context, corID string) context.Context {
	return goutils.ContextWithCorrelationID(ctx, corID)
}

func CorrelationID(ctx context.Context) string {
	return goutils.CorrelationIDFromContext(ctx)
}

func SetLogger(ctx context.Context, l xLog.WrapperLogger) context.Context {
	return xLog.ContextWithLogger(ctx, l)
}

func Logger(ctx context.Context) xLog.WrapperLogger {
	return xLog.LoggerFromContext(ctx)
}
