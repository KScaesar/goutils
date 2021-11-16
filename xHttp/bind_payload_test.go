package xHttp

import (
	"bytes"
	"net/http"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/Min-Feng/goutils/xLog"
	"github.com/Min-Feng/goutils/xTest"
)

func TestBindPayload_Failed(t *testing.T) {
	xLog.SetGlobalLevel("debug")

	gin.SetMode("release")
	router := gin.New()
	router.POST("/hello", RequestIDMiddleware, RecordHttpInfoMiddleware(), bindFailedHandler)

	body := bytes.NewBuffer([]byte(`{"name":"caesar"}`))
	resp, status := xTest.HttpClientDoJson(router, http.MethodPost, "/hello", body)

	expectedResp := `
{
  "code": 1003,
  "msg": "bind payload: Key: 'Person.Age' Error:Field validation for 'Age' failed on the 'required' tag: invalid parameter",
  "data": {}
}`
	assert.JSONEq(t, expectedResp, resp)
	assert.Equal(t, http.StatusBadRequest, status)
}

func bindFailedHandler(c *gin.Context) {
	type Person struct {
		Name string `json:"name"`
		Age  int    `json:"age" binding:"required"`
	}
	payload := new(Person)
	if !BindPayload(c, payload) {
		return
	}

	time.Sleep(281 * time.Millisecond)
	resp := "hello " + payload.Name
	c.JSON(http.StatusOK, NewNormalResponse(resp))
}
