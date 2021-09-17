package xHttp_test

import (
	"os"
	"testing"

	"github.com/Min-Feng/goutils/logger"
)

func TestMain(m *testing.M) {
	logger.TestingMode()

	code := m.Run()
	os.Exit(code)
}
