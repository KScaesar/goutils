package xHttp

import (
	"context"

	"github.com/gin-gonic/gin"
)

const stdCtx = "stdCtx"

// SetStdContextToGin 因為將 std context 放回去 http.Request 的成本太高,
// ginCtx.Request = ginCtx.Request.Clone(stdCtx),
// 所以利用 gin.Context 來傳遞標準庫的 context.Context
//
// 後來發現好像可以用 http.Request.WithContext 來達成目標
// 此函數應該不需要了
func SetStdContextToGin(c *gin.Context, std context.Context) {
	c.Set(stdCtx, std)
}

// GetStdContextFromGin
// 如果在 gin.Context 沒有發現之前塞入的 std context, 則從 Request 取得
func GetStdContextFromGin(c *gin.Context) context.Context {
	v, exists := c.Get(stdCtx)
	if !exists {
		return c.Request.Context()
	}
	return v.(context.Context)
}
