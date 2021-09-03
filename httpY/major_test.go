package httpY_test

import (
	"os"
	"testing"

	"github.com/Min-Feng/goutils/logY"
)

func TestMain(m *testing.M) {
	logY.TestingMode()

	code := m.Run()
	os.Exit(code)
}
