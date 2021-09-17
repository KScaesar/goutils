package xTest

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFakeTimeNow(t *testing.T) {
	now := FakeTimeNow("2021-09-21")
	year, month, day := now().Date()

	assert.Equal(t, 2021, year)
	assert.Equal(t, time.Month(9), month)
	assert.Equal(t, 21, day)
}
