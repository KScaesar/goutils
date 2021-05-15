package logY

import (
	"io"
	"time"

	"github.com/rs/zerolog"
)

func WrapPrototype(prototype zerolog.Logger) WrapperLogger {
	return WrapperLogger{prototype}
}

func New(w io.Writer) WrapperLogger {
	return WrapperLogger{Logger: zerolog.New(w).With().Timestamp().Stack().Logger()}
}

type WrapperLogger struct {
	zerolog.Logger
}

func (l WrapperLogger) Prototype() *zerolog.Logger {
	return &l.Logger
}

// The following, for various scenario

func (l WrapperLogger) ReqBody(body string) WrapperLogger {
	if body == "" {
		body = "empty"
	}
	l.Logger = l.Logger.With().Str("req_body", body).Logger()
	return l
}

func (l WrapperLogger) RespBody(body string) WrapperLogger {
	if body == "" {
		body = "empty"
	}
	l.Logger = l.Logger.With().Str("resp_body", body).Logger()
	return l
}

func (l WrapperLogger) Referrer(refer string) WrapperLogger {
	if refer == "" {
		refer = "empty"
	}
	l.Logger = l.Logger.With().Str("referrer", refer).Logger()
	return l
}

func (l WrapperLogger) ClientIP(ip string) WrapperLogger {
	l.Logger = l.Logger.With().Str("client_ip", ip).Logger()
	return l
}

func (l WrapperLogger) HTTPMethod(method string) WrapperLogger {
	l.Logger = l.Logger.With().Str("http_method", method).Logger()
	return l
}

func (l WrapperLogger) HTTPStatus(status int) WrapperLogger {
	l.Logger = l.Logger.With().Int("http_status", status).Logger()
	return l
}

func (l WrapperLogger) URL(url string) WrapperLogger {
	l.Logger = l.Logger.With().Str("url", url).Logger()
	return l
}

func (l WrapperLogger) CostTime(d time.Duration) WrapperLogger {
	l.Logger = l.Logger.With().Str("cost", d.Truncate(time.Millisecond).String()).Logger()
	return l
}

func (l WrapperLogger) Kind(kind Kind) WrapperLogger {
	l.Logger = l.Logger.With().Str("kind", kind.String()).Logger()
	return l
}

func (l WrapperLogger) TraceID(traceID string) WrapperLogger {
	l.Logger = l.Logger.With().Str("trace_id", traceID).Logger()
	return l
}
