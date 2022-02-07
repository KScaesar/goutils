package goutils

import (
	"io"
	"math/rand"
	"sync"
	"time"

	"github.com/oklog/ulid/v2"
	uuid "github.com/satori/go.uuid"
)

func init() {
	rand.Seed(time.Now().UnixNano())
	entropySeed := time.Now().
		Add(time.Duration(rand.Int63()) * time.Second).
		UnixNano()

	entropyPool.New = func() interface{} {
		return ulid.Monotonic(rand.New(rand.NewSource(entropySeed)), 0)
	}
}

type (
	UUID = string
	ULID = string
)

var entropyPool sync.Pool

func NewULID() ULID {
	entropy := entropyPool.Get().(io.Reader)
	defer entropyPool.Put(entropy)

	now := ulid.Now()
	id := ulid.MustNew(now, entropy)

	return id.String()
}

func NewUUID() UUID {
	return uuid.NewV4().String()
}
