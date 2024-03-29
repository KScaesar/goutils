package xHttp

import (
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/KScaesar/goutils/errors"
	"github.com/KScaesar/goutils/xLog"
	"github.com/KScaesar/goutils/xTest"
)

func TestSendErrorResponse(t *testing.T) {
	xLog.SetGlobalLevel("panic")

	gin.SetMode("release")
	router := gin.New()
	router.POST("/hello", MiddlewareCorrelationID, MiddlewareRecordHttpInfo(), helloHandlerUseCaseFailed)
	resp, status := xTest.HttpClientDoJson(router, http.MethodPost, "/hello", nil)

	expectedResp := `
{
  "code": 1001,
  "msg": "repo: sql statement invalid: system failed",
  "data": {}
}`
	assert.JSONEq(t, expectedResp, resp)
	assert.Equal(t, http.StatusInternalServerError, status)
}

func helloHandlerUseCaseFailed(c *gin.Context) {
	repo := func() error { return errors.Wrap(errors.ErrSystem, "sql statement invalid") }

	uc := func() error {
		err := repo()
		if err != nil {
			return errors.WrapMessage(err, "repo")
		}
		return nil
	}

	if err := uc(); err != nil {
		SendErrorResponse(c, err)
		return // 一定要記得 return
	}
}
