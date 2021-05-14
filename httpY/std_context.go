package httpY

import (
	"context"

	"github.com/gin-gonic/gin"
)

const stdCtx = "stdCtx"

// SetStdContext 因為將 std context 放回去 http.Request 的成本太高,
// ginCtx.Request = ginCtx.Request.Clone(stdCtx),
// 所以利用 gin.Context 來傳遞標準庫的 context.Context
func SetStdContext(c *gin.Context, std context.Context) {
	c.Set(stdCtx, std)
}

// GetStdContext 因為將 std context 放回去 http.Request 的成本太高,
// ginCtx.Request = ginCtx.Request.Clone(stdCtx),
// 所以利用 gin.Context 來傳遞標準庫的 context.Context,
// 如果在 gin.Context 沒有發現 std context, 則從 Request 取得
func GetStdContext(c *gin.Context) context.Context {
	v, exists := c.Get(stdCtx)
	if !exists {
		return c.Request.Context()
	}
	return v.(context.Context)
}
