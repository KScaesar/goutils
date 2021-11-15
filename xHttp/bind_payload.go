package xHttp

import (
	"github.com/gin-gonic/gin"

	"github.com/Min-Feng/goutils/errors"
)

// BindPayload
// 若 client 端, 沒有設置 application/json,
// 不僅 ShouldBind 抓不到資料 且 error == nil
func BindPayload(c *gin.Context, obj interface{}) bool {
	if err := c.ShouldBind(obj); err != nil {
		Err := errors.Wrap(errors.ErrInvalidParams, "bind payload: %v", err)
		sendErrorResponseBase(c, Err, 1)
		return false
	}
	return true
}
