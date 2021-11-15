package xLog

import (
	"context"

	"github.com/google/uuid"
)

type requestKey struct{}

func ContextWithRequestID(ctx context.Context, reqID string) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	return context.WithValue(ctx, requestKey{}, reqID)
}

func RequestIDFromContext(ctx context.Context) string {
	reqID, ok := ctx.Value(requestKey{}).(string)
	if ok {
		return reqID
	}
	return ""
}

func NewRequestID() string {
	return uuid.NewString()
}
