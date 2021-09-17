package xHttp

import (
	"github.com/gin-gonic/gin"

	"github.com/Min-Feng/goutils/errors"
	"github.com/Min-Feng/goutils/xLog"
)

func SendErrorResponse(c *gin.Context, err error) {
	SendErrorResponseBase(c, err, 1, xLog.KindApplication)
}

func SendErrorResponseBase(c *gin.Context, err error, skip int, kind xLog.Kind) {
	if err == nil {
		return
	}

	xLog.LoggerFromContext(c.Request.Context()).
		Kind(kind).
		Unwrap().Err(err).Caller(skip + 1).Send()

	c.JSON(errors.HTTPStatus(err), NewErrorResponse(err))
	c.Abort()
}
