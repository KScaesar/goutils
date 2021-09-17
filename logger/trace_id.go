package logger

import (
	"context"
	_ "unsafe"

	"github.com/google/uuid"
)

type TraceKey struct{}

func NewTraceIDCtx(ctx context.Context, traceID string) (traceIDCtx context.Context) {
	if ctx == nil {
		ctx = context.Background()
	}
	return context.WithValue(ctx, TraceKey{}, traceID)
}

func TraceIDFromCtx(traceIDCtx context.Context) string {
	traceID, ok := traceIDCtx.Value(TraceKey{}).(string)
	if ok {
		return traceID
	}
	return ""
}

func NewTraceID() string {
	return uuid.NewString()
}
