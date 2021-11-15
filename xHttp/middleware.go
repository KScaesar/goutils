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

const RequestIDHeaderKey = "X-RequestID"

func RequestIDMiddleware(c *gin.Context) {
	reqID := RequestIDFromHeader(c.Request)
	c.Writer.Header().Add(RequestIDHeaderKey, reqID)

	ctx := xLog.ContextWithRequestID(c.Request.Context(), reqID)
	logger := xLog.Logger().RequestID(reqID)
	logCtx := xLog.ContextWithLogger(ctx, logger)
	c.Request = c.Request.WithContext(logCtx)

	c.Next()
}

func RequestIDFromHeader(r *http.Request) (reqID string) {
	reqID = r.Header.Get(RequestIDHeaderKey)
	if reqID == "" {
		reqID = xLog.NewRequestID()
	}
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

func RecordHttpInfoMiddleware() gin.HandlerFunc {
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

		var reqBody bytes.Buffer
		var respWriter respMultiWriter
		if xLog.IsDebugLevel() {
			teeReader := io.TeeReader(c.Request.Body, &reqBody)
			c.Request.Body = ioutil.NopCloser(teeReader)

			respWriter.ResponseWriter = c.Writer
			c.Writer = &respWriter
		}

		c.Next()

		logger := xLog.LoggerFromContext(c.Request.Context())

		info := &xLog.HttpMetricInfo{
			Method:   c.Request.Method,
			URL:      c.Request.URL.Redacted(),
			ClientIP: c.ClientIP(),
			Referrer: c.Request.Referer(),
			Status:   c.Writer.Status(),
			TimeCost: time.Now().Sub(start),
		}

		if xLog.IsDebugLevel() && !containKeyword(reqBody.Bytes()) {
			debug := &xLog.HttpMetricDebug{
				ReqBody:  reqBody.String(),
				RespBody: respWriter.body.String(),
			}
			logger.RecordHttpForDebug(info, debug).Unwrap().Debug().Send()
		}

		if c.Writer.Status() >= http.StatusBadRequest {
			logger.RecordHttp(info).Unwrap().Error().Send()
			return
		}

		logger.RecordHttp(info).Unwrap().Info().Send()
	}
}
