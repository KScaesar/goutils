package httpY

import (
	"github.com/Min-Feng/goutils/errorY"
	"github.com/Min-Feng/goutils/logY"

	"github.com/gin-gonic/gin"
)

// BindPayload
// 若 client 端, 沒有設置 application/json,
// 不僅 ShouldBind 抓不到資料 且 error == nil
func BindPayload(c *gin.Context, obj interface{}) bool {
	if err := c.ShouldBind(obj); err != nil {
		Err := errorY.Wrap(errorY.ErrInvalidParams, "bind payload: %v", err)

		logY.FromCtx(GetStdContext(c)).
			Kind(logY.KindHTTP).Prototype().Err(Err).Caller(1).Send()

		c.JSON(errorY.HTTPStatus(Err), NewErrorResponse(Err))
		c.Abort()
		return false
	}
	return true
}
