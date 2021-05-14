package httpY

import (
	"bytes"
	"net/http"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/Min-Feng/goutils/logY"
	"github.com/Min-Feng/goutils/testingY"
)

func TestMiddleware_BindPayloadFailed(t *testing.T) {
	logY.TestingMode()

	gin.SetMode("release")
	router := gin.New()
	router.POST("/hello", TraceIDMiddleware, RecordHTTPInfoMiddleware, ErrorResponseMiddleware, helloHandlerFailed)

	body := bytes.NewBuffer([]byte(`{"name":"caesar"}`))
	resp, status := testingY.HttpClientDoJson(router, http.MethodPost, "/hello", body)

	expectedResp := `
{
  "code": 10002,
  "msg": "bind payload: Key: 'Person.Age' Error:Field validation for 'Age' failed on the 'required' tag: invalid params",
  "data": {}
}`
	assert.JSONEq(t, expectedResp, resp)
	assert.Equal(t, http.StatusBadRequest, status)
}

func helloHandlerFailed(c *gin.Context) {
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

func TestMiddleware_Success(t *testing.T) {
	// logY.FixBugMode()

	gin.SetMode("release")
	router := gin.New()
	router.POST("/hello", TraceIDMiddleware, RecordHTTPInfoMiddleware, ErrorResponseMiddleware, helloHandlerSuccess)

	body := bytes.NewBufferString(`{"name":"caesar"}`)
	resp, status := testingY.HttpClientDoJson(router, http.MethodPost, "/hello", body)

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

	resp := "hello " + payload.Name
	c.JSON(http.StatusOK, NewNormalResponse(resp))
}
