package xHttp

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/Min-Feng/goutils/xLog"
)

const TraceIDHeaderKey = "X-Trace"

func TraceIDFromHeader(r *http.Request, w http.ResponseWriter) (traceID string) {
	traceID = r.Header.Get(TraceIDHeaderKey)
	if traceID == "" {
		traceID = xLog.NewTraceID()
	}
	w.Header().Add(TraceIDHeaderKey, traceID)
	return
}

func TraceIDMiddleware(c *gin.Context) {
	traceID := TraceIDFromHeader(c.Request, c.Writer)
	traceIDCtx := xLog.ContextWithTraceID(c.Request.Context(), traceID)

	log := xLog.Logger().TraceID(traceID)
	logCtx := xLog.ContextWithLogger(traceIDCtx, log)

	c.Request = c.Request.WithContext(logCtx)
	c.Next()
}

type respMultiWriter struct {
	gin.ResponseWriter
	body bytes.Buffer
}

func (w *respMultiWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func RecordHTTPInfoMiddleware() gin.HandlerFunc {
	// keywords 若放在匿名函數 ContainKeyword 裡面, 會造成重複 allocate memory, 利用閉包鎖住變數位址
	keywords := [][]byte{
		[]byte("password"),
	}

	ContainKeyword := func(reqBody []byte) bool {
		for _, keyword := range keywords {
			if bytes.Contains(reqBody, keyword) {
				return true
			}
		}
		return false
	}

	return func(c *gin.Context) {
		start := time.Now()

		var reqBody bytes.Buffer
		var respWriter respMultiWriter
		if xLog.IsDebugLevel() {
			teeReader := io.TeeReader(c.Request.Body, &reqBody)
			c.Request.Body = ioutil.NopCloser(teeReader)

			respWriter.ResponseWriter = c.Writer
			c.Writer = &respWriter
		}

		c.Next()

		log := xLog.LoggerFromContext(c.Request.Context())

		m1 := &xLog.HttpMetricNormal{
			Method:   c.Request.Method,
			URL:      c.Request.URL.Redacted(),
			ClientIP: c.ClientIP(),
			Referrer: c.Request.Referer(),
			Status:   c.Writer.Status(),
			TimeCost: time.Now().Sub(start),
		}

		if xLog.IsDebugLevel() && !ContainKeyword(reqBody.Bytes()) {
			m2 := &xLog.HttpMetricDebug{
				ReqBody:  reqBody.String(),
				RespBody: respWriter.body.String(),
			}
			log.RecordHttpInfo(m1, m2).Unwrap().Debug().Send()
		}

		status := c.Writer.Status()
		if status >= http.StatusBadRequest {
			log.RecordHttpInfo(m1, nil).Unwrap().Error().Send()
			return
		}

		log.RecordHttpInfo(m1, nil).Unwrap().Info().Send()
	}

}
