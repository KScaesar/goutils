package xHttp

import (
	"bytes"

	"github.com/Min-Feng/goutils/errors"
)

func init() {
	errors.RegisterFrameFilter(errorFilterGinNext())
}

func errorFilterGinNext() errors.FrameFilter {
	target := []byte("gin.(*Context).Next")
	return func(frame []byte) bool {
		return bytes.Contains(frame, target)
	}
}
