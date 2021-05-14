package logY

import (
	"bytes"
	"net/http"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"

	"github.com/Min-Feng/goutils/errorY"
	"github.com/Min-Feng/goutils/testingY"
)

//go:generate go test -trimpath -run=Error -v github.com/Min-Feng/goutils/logY
func TestError(t *testing.T) {
	writer := &bytes.Buffer{}
	_default = New(writer)
	defer DefaultMode()

	zerolog.TimeFieldFormat = CustomTimeFormat // DefaultMode() show=human, 影響 zerolog.TimeFieldFormat
	zerolog.TimestampFunc = testingY.FakeTimeNow("2021-12-14")

	err := errorY.New(8787, http.StatusInternalServerError, "json failed")
	Err := errorY.Wrap(err, "unit test")
	Error(Err).Send()

	expected := `
{
  "level": "error",
  "err_code": 8787,
  "error": "unit test: json failed",
  "caller": "github.com/Min-Feng/goutils/logY/default_test.go:26",
  "timestamp": "2021-12-14 00:00:00+08:00",
  "stack": [
    [
      "github.com/Min-Feng/goutils/logY.TestError github.com/Min-Feng/goutils/logY/default_test.go:25 ",
      "testing.tRunner testing/testing.go:1123 ",
      "runtime.goexit runtime/asm_amd64.s:1374 "
    ]
  ]
}`

	assert.JSONEq(t, expected, writer.String())
}
