package xLog

import (
	"io"
	"time"

	"github.com/rs/zerolog"
)

func WrapPrototype(prototype *zerolog.Logger) WrapperLogger {
	return WrapperLogger{*prototype}
}

func NewLogger(w io.Writer) WrapperLogger {
	return WrapperLogger{
		prototype: zerolog.New(w).With().Timestamp().Stack().Logger(),
	}
}

type WrapperLogger struct {
	prototype zerolog.Logger
}

func (l WrapperLogger) Unwrap() *zerolog.Logger {
	return &l.prototype
}

func (l WrapperLogger) ReqBody(body string) WrapperLogger {
	if body == "" {
		body = "empty"
	}
	l.prototype = l.prototype.With().Str("req_body", body).Logger()
	return l
}

func (l WrapperLogger) RespBody(body string) WrapperLogger {
	if body == "" {
		body = "empty"
	}
	l.prototype = l.prototype.With().Str("resp_body", body).Logger()
	return l
}

func (l WrapperLogger) Referrer(refer string) WrapperLogger {
	if refer == "" {
		refer = "empty"
	}
	l.prototype = l.prototype.With().Str("referrer", refer).Logger()
	return l
}

func (l WrapperLogger) ClientIP(ip string) WrapperLogger {
	l.prototype = l.prototype.With().Str("client_ip", ip).Logger()
	return l
}

func (l WrapperLogger) HttpMethod(method string) WrapperLogger {
	l.prototype = l.prototype.With().Str("http_method", method).Logger()
	return l
}

func (l WrapperLogger) HttpStatus(status int) WrapperLogger {
	l.prototype = l.prototype.With().Int("http_status", status).Logger()
	return l
}

func (l WrapperLogger) URL(url string) WrapperLogger {
	l.prototype = l.prototype.With().Str("url", url).Logger()
	return l
}

func (l WrapperLogger) CostTime(d time.Duration) WrapperLogger {
	l.prototype = l.prototype.With().Str("cost", d.Truncate(time.Millisecond).String()).Logger()
	return l
}

func (l WrapperLogger) TriggerKind(kind TriggerKind) WrapperLogger {
	l.prototype = l.prototype.With().Str("trigger_kind", kind.String()).Logger()
	return l
}

func (l WrapperLogger) RequestID(requestID string) WrapperLogger {
	l.prototype = l.prototype.With().Str("request_id", requestID).Logger()
	return l
}
