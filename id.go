package goutils

import (
	"io"
	"math/rand"
	"sync"
	"time"

	"github.com/oklog/ulid/v2"
	uuid "github.com/satori/go.uuid"
)

type (
	UUID = string
	ULID = string
)

var defaultEntropyPool = newULIDEntropyPool()

func NewULID() ULID {
	entropy := defaultEntropyPool.Get().(io.Reader)
	defer defaultEntropyPool.Put(entropy)

	now := ulid.Now()
	id := ulid.MustNew(now, entropy)

	return id.String()
}

func NewULIDForReplay(t time.Time) ULID {
	entropy := newULIDEntropyFactory(time.Now)()
	id := ulid.MustNew(ulid.Timestamp(t), entropy)
	return id.String()
}

func newULIDEntropyPool() sync.Pool {
	factory := newULIDEntropyFactory(time.Now)
	return sync.Pool{
		New: func() interface{} {
			return factory()
		},
	}
}

func newULIDEntropyFactory(timeNowFn func() time.Time) func() *ulid.MonotonicEntropy {
	rand.Seed(timeNowFn().UnixNano())
	randomTime := time.Duration(rand.Int63()) * time.Millisecond
	seed := timeNowFn().Add(randomTime).UnixNano()

	return func() *ulid.MonotonicEntropy {
		seed++
		return ulid.Monotonic(rand.New(rand.NewSource(seed)), 0)
	}
}

func NewUUID() UUID {
	return uuid.NewV4().String()
}
