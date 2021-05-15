package httpY

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/Min-Feng/goutils"
	"github.com/Min-Feng/goutils/errorY"
	"github.com/Min-Feng/goutils/logY"
)

func TraceIDMiddleware(c *gin.Context) {
	ctx := c.Request.Context()

	traceID := c.GetHeader("X-TraceID")
	if traceID == "" {
		traceID = goutils.NewTraceID()
	}
	traceIDCtx := goutils.TraceIDWithCtx(ctx, traceID)

	SetStdContext(c, logY.Logger().TraceID(traceID).WithCtx(traceIDCtx))
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

	if logY.IsDebugLevel() {
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

func ErrorResponseMiddleware(c *gin.Context) {
	c.Next()

	if len(c.Errors) == 0 {
		return
	}
	err := c.Errors[0].Err

	log := logY.FromCtx(GetStdContext(c)).Kind(logY.KindApplication)

	log.Err(err).Send()
	c.JSON(errorY.HTTPStatus(err), NewErrorResponse(err))

	if len(c.Errors) > 1 {
		for i, ginErr := range c.Errors {
			Err := errorY.Wrap(errorY.ErrSystem, "not should have many error: [%d] %v", i, ginErr)
			log.Err(Err).Send()
		}
	}
}
