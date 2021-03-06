package xHttp

import (
	"bytes"
	"net/http"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/KScaesar/goutils/xLog"
	"github.com/KScaesar/goutils/xTest"
)

func TestMiddleware_Success(t *testing.T) {
	xLog.SetGlobalLevel("panic")

	gin.SetMode("release")
	router := gin.New()
	router.POST("/hello", MiddlewareCorrelationID, MiddlewareRecordHttpInfo(), helloHandlerSuccess)

	body := bytes.NewBufferString(`{"name":"caesar"}`)
	resp, status := xTest.HttpClientDoJson(router, http.MethodPost, "/hello", body)

	expectedResp := `
{
  "code": 0,
  "msg": "ok",
  "data": "hello caesar"
}`
	assert.JSONEq(t, expectedResp, resp)
	assert.Equal(t, http.StatusOK, status)
}

func helloHandlerSuccess(c *gin.Context) {
	type Person struct {
		Name string `json:"name"`
	}
	payload := new(Person)
	if !BindPayload(c, payload) {
		return
	}

	time.Sleep(234 * time.Millisecond)
	resp := "hello " + payload.Name
	c.JSON(http.StatusOK, NewNormalResponse(resp))
}
