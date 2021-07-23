package httpY

import (
	"bytes"
	"context"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/Min-Feng/goutils/logY"
)

func TraceIDMiddleware(c *gin.Context) {
	traceIDCtx, traceID := NewTraceIDCtx(c.Request, c.Writer)
	logger := logY.Logger().TraceID(traceID)
	c.Request = c.Request.WithContext(logY.NewLogContext(traceIDCtx, logger))
	c.Next()
}

const TraceIDHeaderKey = "X-Trace"

func NewTraceIDCtx(r *http.Request, w http.ResponseWriter) (traceIDCtx context.Context, traceID string) {
	traceID = r.Header.Get(TraceIDHeaderKey)
	if traceID == "" {
		traceID = logY.NewTraceID()
	}

	w.Header().Add(TraceIDHeaderKey, traceID)
	traceIDCtx = logY.NewTraceIDCtx(r.Context(), traceID)
	return
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
		if logY.IsDebugLevel() {
			teeReader := io.TeeReader(c.Request.Body, &reqBody)
			c.Request.Body = ioutil.NopCloser(teeReader)

			respWriter.ResponseWriter = c.Writer
			c.Writer = &respWriter
		}

		c.Next()

		log := logY.FromCtx(c.Request.Context())

		m1 := &logY.HttpMetricNormal{
			Method:   c.Request.Method,
			URL:      c.Request.URL.Redacted(),
			ClientIP: c.ClientIP(),
			Referrer: c.Request.Referer(),
			Status:   c.Writer.Status(),
			TimeCost: time.Now().Sub(start),
		}

		if logY.IsDebugLevel() && !ContainKeyword(reqBody.Bytes()) {
			m2 := &logY.HttpMetricDebug{
				ReqBody:  reqBody.String(),
				RespBody: respWriter.body.String(),
			}
			log.RecordHttpInfo(m1, m2).Prototype().Debug().Send()
		}

		status := c.Writer.Status()
		if status >= http.StatusBadRequest {
			log.RecordHttpInfo(m1, nil).Prototype().Error().Send()
			return
		}

		log.RecordHttpInfo(m1, nil).Prototype().Info().Send()
	}

}
