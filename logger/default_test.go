package logger

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

//go:generate go test -trimpath -run=TestError -v github.com/Min-Feng/goutils/logY
func TestError(t *testing.T) {
	writer := &bytes.Buffer{}
	default_ = New(writer)
	defer DefaultMode()

	zerolog.TimeFieldFormat = CustomTimeFormat // DefaultMode() show=human, 影響 zerolog.TimeFieldFormat
	zerolog.TimestampFunc = xTest.FakeTimeNow("2021-12-14")

	expected := `
{
  "level": "error",
  "caller": "github.com/Min-Feng/goutils/logY/default_test.go:44",
  "timestamp": "2021-12-14 00:00:00+08:00",
  "error": {
    "msg": "unit test: json failed",
    "code": 8787
  },
  "stack": [
    [
      "github.com/Min-Feng/goutils/logY.TestError github.com/Min-Feng/goutils/logY/default_test.go:44 ",
      "testing.tRunner testing/testing.go:1193 ",
      "runtime.goexit runtime/asm_amd64.s:1371 "
    ]
  ]
}`

	err := errors.New(8787, http.StatusInternalServerError, "json failed")
	Err(errors.Wrap(err, "unit test")).Send()
	actual := writer.String()

	errPath := "error"
	assert.JSONEq(t, gjson.Get(expected, errPath).Raw, gjson.Get(actual, errPath).Raw)
	stackPath := "stack.0.0"
	assert.Equal(t, gjson.Get(expected, stackPath).Raw, gjson.Get(actual, stackPath).Raw)
}
