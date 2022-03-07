package xHttp

import (
	"github.com/gin-gonic/gin"

	"github.com/KScaesar/goutils/errors"
	"github.com/KScaesar/goutils/xLog"
)

func SendErrorResponse(c *gin.Context, err error) {
	sendErrorResponseBase(c, err, 1)
}

func sendErrorResponseBase(c *gin.Context, err error, skip int) {
	if err == nil {
		return
	}

	xLog.LoggerFromContext(c.Request.Context()).
		Unwrap().
		Err(err).
		Caller(skip + 1).
		Send()

	c.JSON(errors.HttpStatus(err), NewErrorResponse(err))
	c.Abort()
}
