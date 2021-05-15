package goutils

import "context"

type traceKey struct{}

func TraceIDFromCtx(ctx context.Context) string {
	traceID, ok := ctx.Value(traceKey{}).(string)
	if ok {
		return traceID
	}
	return "empty"
}

func TraceIDWithCtx(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, traceKey{}, traceID)
}

func NewTraceID() string {
	return "xxxccc1122"
}
