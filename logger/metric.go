package logger

import "time"

type HttpMetricNormal struct {
	Method   string
	URL      string
	ClientIP string
	Referrer string
	Status   int
	TimeCost time.Duration
}

type HttpMetricDebug struct {
	ReqBody  string
	RespBody string
}

func (l WrapperLogger) RecordHttpInfo(normal *HttpMetricNormal, debug *HttpMetricDebug) WrapperLogger {
	log := l

	if normal != nil {
		log = log.
			Kind(KindHTTP).
			HTTPMethod(normal.Method).
			URL(normal.URL).
			ClientIP(normal.ClientIP).
			Referrer(normal.Referrer).
			HTTPStatus(normal.Status).
			CostTime(normal.TimeCost)
	}

	if debug != nil {
		log = log.
			ReqBody(debug.ReqBody).
			RespBody(debug.RespBody)
	}

	return log
}
