package goutils

import (
	"testing"
	"time"

	"github.com/oklog/ulid/v2"
	"github.com/stretchr/testify/assert"
)

func Test_newULIDEntropyFactory_GetDifferentULID_ByDifferentEntropy(t *testing.T) {
	now := ulid.Now()
	entropyFactory := newULIDEntropyFactory(time.Now)
	qty := 5
	result := make(map[string]bool, qty)

	for i := 0; i < qty; i++ {
		entropy := entropyFactory()
		id := ulid.MustNew(now, entropy).String()
		result[id] = true
	}

	actualQty := len(result)
	assert.Equal(t, qty, actualQty)
}

func Test_newULIDEntropyFactory_GetDifferentULID_BySameEntropy(t *testing.T) {
	now := ulid.Now()
	entropyFactory := newULIDEntropyFactory(time.Now)
	entropy := entropyFactory()
	qty := 5
	result := make(map[string]bool, qty)

	for i := 0; i < qty; i++ {
		id := ulid.MustNew(now, entropy).String()
		result[id] = true
	}

	actualQty := len(result)
	assert.Equal(t, qty, actualQty)
}
