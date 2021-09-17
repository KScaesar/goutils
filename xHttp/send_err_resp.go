package xHttp

import (
	"github.com/gin-gonic/gin"

	"github.com/Min-Feng/goutils/errors"
	"github.com/Min-Feng/goutils/logger"
)

func SendErrorResponse(c *gin.Context, err error) {
	SendErrorResponseBase(c, err, 1, logger.KindApplication)
}

func SendErrorResponseBase(c *gin.Context, err error, skip int, kind logger.Kind) {
	if err == nil {
		return
	}

	logger.FromCtx(c.Request.Context()).
		Kind(kind).
		Prototype().Err(err).Caller(skip + 1).Send()

	c.JSON(errors.HTTPStatus(err), NewErrorResponse(err))
	c.Abort()
}
