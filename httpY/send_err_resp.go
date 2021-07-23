package httpY

import (
	"github.com/gin-gonic/gin"

	"github.com/Min-Feng/goutils/errorY"
	"github.com/Min-Feng/goutils/logY"
)

func SendErrorResponse(c *gin.Context, err error) {
	SendErrorResponseBase(c, err, 1, logY.KindApplication)
}

func SendErrorResponseBase(c *gin.Context, err error, skip int, kind logY.Kind) {
	if err == nil {
		return
	}

	logY.FromCtx(c.Request.Context()).
		Kind(kind).
		Prototype().Err(err).Caller(skip + 1).Send()

	c.JSON(errorY.HTTPStatus(err), NewErrorResponse(err))
	c.Abort()
}
