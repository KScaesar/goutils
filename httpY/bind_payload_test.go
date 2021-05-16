package httpY

import (
	"bytes"
	"net/http"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/Min-Feng/goutils/testingY"
)

func TestBindPayload_Failed(t *testing.T) {
	// logY.FixBugMode()

	gin.SetMode("release")
	router := gin.New()
	router.POST("/hello", TraceIDMiddleware, RecordHTTPInfoMiddleware, bindFailedHandler)

	body := bytes.NewBuffer([]byte(`{"name":"caesar"}`))
	resp, status := testingY.HttpClientDoJson(router, http.MethodPost, "/hello", body)

	expectedResp := `
{
  "code": 1002,
  "msg": "bind payload: Key: 'Person.Age' Error:Field validation for 'Age' failed on the 'required' tag: invalid params",
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
