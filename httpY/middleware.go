package httpY

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/Min-Feng/goutils/logY"
)

func TraceIDMiddleware(c *gin.Context) {
	traceID, ctx := logY.TraceIDFromHTTP(c.Request, c.Writer)
	SetStdContext(
		c,
		logY.Logger().TraceID(traceID).WithCtx(ctx),
	)
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

func RecordHTTPInfoMiddleware(c *gin.Context) {
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

	log := logY.FromCtx(GetStdContext(c))

	m1 := &logY.HttpMetricNormal{
		Method:   c.Request.Method,
		URL:      c.Request.URL.Redacted(),
		ClientIP: c.ClientIP(),
		Referrer: c.Request.Referer(),
		Status:   c.Writer.Status(),
		TimeCost: time.Now().Sub(start),
	}

	if logY.IsDebugLevel() && !bytes.Contains(reqBody.Bytes(), []byte("password")) {
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
