package xLog

import (
	"bytes"
	"net/http"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/tidwall/gjson"

	"github.com/KScaesar/goutils/errors"
	"github.com/KScaesar/goutils/xTest"
)

//go:generate go test -trimpath -run=TestError -v github.com/KScaesar/goutils/xLog
func TestError(t *testing.T) {
	writer := &bytes.Buffer{}
	Init("debug", false)
	SetDefaultLogger(NewLogger(writer))

	zerolog.TimestampFunc = xTest.FakeTimeNow("2021-12-14")

	expected := `
{
  "level": "error",
  "stack": [
    [
      "github.com/KScaesar/goutils/xLog.TestError github.com/KScaesar/goutils/xLog/default_test.go:43 ",
      "testing.tRunner testing/testing.go:1193 ",
      "runtime.goexit runtime/asm_amd64.s:1371 "
    ]
  ],
  "my_error": {
    "err_msg": "unit test: json failed",
    "err_code": 8787
  },
  "caller": "github.com/KScaesar/goutils/xLog/default_test.go:42",
  "timestamp": 1639440000000
}`

	err := errors.New(8787, http.StatusInternalServerError, "json failed")
	Err(errors.Wrap(err, "unit test")).Send()
	actual := writer.String()

	errPath := "my_error"
	assert.JSONEq(t, gjson.Get(expected, errPath).Raw, gjson.Get(actual, errPath).Raw)
	stackPath := "stack.0.0"
	assert.Equal(t, gjson.Get(expected, stackPath).Raw, gjson.Get(actual, stackPath).Raw)
}
