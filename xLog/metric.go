package xLog

import "time"

type HttpMetricInfo struct {
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

func (l WrapperLogger) RecordHttpForDebug(debug *HttpMetricDebug) WrapperLogger {
	return l.
		ReqBody(debug.ReqBody).
		RespBody(debug.RespBody)
}

func (l WrapperLogger) RecordHttp(info *HttpMetricInfo) WrapperLogger {
	return l.
		HttpMethod(info.Method).
		URL(info.URL).
		ClientIP(info.ClientIP).
		Referrer(info.Referrer).
		HttpStatus(info.Status).
		CostTime(info.TimeCost)
}
