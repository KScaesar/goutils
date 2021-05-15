package httpY

import (
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/Min-Feng/goutils/errorY"
	"github.com/Min-Feng/goutils/testingY"
)

func TestSendErrorResponse(t *testing.T) {
	// logY.FixBugMode()

	gin.SetMode("release")
	router := gin.New()
	router.POST("/hello", helloHandlerUseCaseFailed)
	resp, status := testingY.HttpClientDoJson(router, http.MethodPost, "/hello", nil)

	expectedResp := `
{
  "code": 10001,
  "msg": "sql statement invalid: system failed",
  "data": {}
}`
	assert.JSONEq(t, expectedResp, resp)
	assert.Equal(t, http.StatusInternalServerError, status)
}

func helloHandlerUseCaseFailed(c *gin.Context) {
	repo := func() error { return errorY.Wrap(errorY.ErrSystem, "sql statement invalid") }

	uc := func() error {
		err := repo()
		if err != nil {
			return errorY.WrapMessage(err, "repo")
		}
		return nil
	}

	if err := uc(); err != nil {
		SendErrorResponse(c, err)
		return // 一定要記得 return
	}
}
