package xLog

import (
	"bytes"
	"net/http"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/tidwall/gjson"

	"github.com/Min-Feng/goutils/errors"
	"github.com/Min-Feng/goutils/xTest"
)

//go:generate go test -trimpath -run=TestError -v github.com/Min-Feng/goutils/xLog
func TestError(t *testing.T) {
	writer := &bytes.Buffer{}
	Init("error", false)
	SetDefaultLogger(NewLogger(writer))

	zerolog.TimestampFunc = xTest.FakeTimeNow("2021-12-14")

	expected := `
{
  "level": "error",
  "stack": [
    [
      "github.com/Min-Feng/goutils/xLog.TestError github.com/Min-Feng/goutils/xLog/default_test.go:43 ",
      "testing.tRunner testing/testing.go:1193 ",
      "runtime.goexit runtime/asm_amd64.s:1371 "
    ]
  ],
  "error": {
    "err_msg": "unit test: json failed",
    "err_code": 8787
  },
  "caller": "github.com/Min-Feng/goutils/xLog/default_test.go:42",
  "timestamp": 1639440000000
}`

	err := errors.New(8787, http.StatusInternalServerError, "json failed")
	Err(errors.Wrap(err, "unit test")).Send()
	actual := writer.String()

	errPath := "error"
	assert.JSONEq(t, gjson.Get(expected, errPath).Raw, gjson.Get(actual, errPath).Raw)
	stackPath := "stack.0.0"
	assert.Equal(t, gjson.Get(expected, stackPath).Raw, gjson.Get(actual, stackPath).Raw)
}
