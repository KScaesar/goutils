package goutils

import "context"

func NewCorrelationID() string {
	return NewUUID()
}

type correlationIDKey struct{}

func ContextWithCorrelationID(ctx context.Context, corID string) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	if corID == "" {
		corID = "param_empty"
	}
	return context.WithValue(ctx, correlationIDKey{}, corID)
}

func CorrelationIDFromContext(ctx context.Context) string {
	corID, ok := ctx.Value(correlationIDKey{}).(string)
	if ok {
		return corID
	}
	return "no_assign"
}
