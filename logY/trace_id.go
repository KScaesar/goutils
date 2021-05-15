package logY

import (
	"context"
	"net/http"
	_ "unsafe"

	"github.com/google/uuid"
)

var TraceKey = "X-Trace"

func TraceIDFromHTTP(r *http.Request, w http.ResponseWriter) (traceID string, ctx context.Context) {
	traceID = r.Header.Get(TraceKey)
	if traceID == "" {
		traceID = NewTraceID()
		w.Header().Add(TraceKey, traceID)
	}
	ctx = context.WithValue(r.Context(), &TraceKey, traceID)
	return
}

func TraceIDFromCtx(ctx context.Context) string {
	traceID, ok := ctx.Value(&TraceKey).(string)
	if ok {
		return traceID
	}
	return ""
}

func NewTraceID() string {
	return uuid.NewString()
}
