package xHttp

import (
	"bytes"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/KScaesar/goutils"
	"github.com/KScaesar/goutils/xLog"
)

const CorrelationIDHeaderKey = "X-CorrelationID"

func MiddlewareCorrelationID(c *gin.Context) {
	corID := c.Request.Header.Get(CorrelationIDHeaderKey)
	if corID == "" {
		corID = goutils.NewCorrelationID()
	}
	c.Writer.Header().Add(CorrelationIDHeaderKey, corID)

	logger := xLog.Logger().
		CorrelationID(corID)

	ctx := c.Request.Context()
	ctx1 := goutils.ContextWithCorrelationID(ctx, corID)
	ctx2 := xLog.ContextWithLogger(ctx1, logger)
	c.Request = c.Request.WithContext(ctx2)

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

type requestMultiReader struct {
	io.ReadCloser
	body bytes.Buffer
}

func (r *requestMultiReader) Read(p []byte) (n int, err error) {
	return io.TeeReader(r.ReadCloser, &r.body).Read(p)
}

func (r *requestMultiReader) Close() error {
	return r.ReadCloser.Close()
}

func MiddlewareRecordHttpInfo() gin.HandlerFunc {
	// keywords 若放在匿名函數 containKeyword 裡面, 會造成重複 allocate memory, 利用閉包鎖住變數位址
	keywords := [][]byte{
		[]byte("password"),
	}

	containKeyword := func(reqBody []byte) bool {
		for _, keyword := range keywords {
			if bytes.Contains(reqBody, keyword) {
				return true
			}
		}
		return false
	}

	return func(c *gin.Context) {
		start := time.Now()

		var reqReader requestMultiReader
		var respWriter respMultiWriter
		if xLog.IsDebugLevel() {
			reqReader.ReadCloser = c.Request.Body
			c.Request.Body = &reqReader

			respWriter.ResponseWriter = c.Writer
			c.Writer = &respWriter
		}

		c.Next()

		logger := xLog.LoggerFromContext(c.Request.Context()).TriggerKind(xLog.TriggerKindHttp)

		info := &xLog.HttpMetricInfo{
			Method:   c.Request.Method,
			URL:      c.Request.URL.Redacted(),
			ClientIP: c.ClientIP(),
			Referrer: c.Request.Referer(),
			Status:   c.Writer.Status(),
			TimeCost: time.Now().Sub(start),
		}
		logger = logger.RecordHttp(info)

		if xLog.IsDebugLevel() && !containKeyword(reqReader.body.Bytes()) {
			debug := &xLog.HttpMetricDebug{
				ReqBody:  reqReader.body.String(),
				RespBody: respWriter.body.String(),
			}
			logger = logger.RecordHttpForDebug(debug)
		}

		if c.Writer.Status() >= http.StatusBadRequest {
			logger.RecordHttp(info).Unwrap().Error().Send()
			return
		}

		logger.RecordHttp(info).Unwrap().Info().Send()
	}
}
