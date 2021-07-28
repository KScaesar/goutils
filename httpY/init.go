package httpY

import (
	"bytes"

	"github.com/Min-Feng/goutils/errorY"
)

func init() {
	errorY.RegisterFrameFilter(ginNextFilter())
}

func ginNextFilter() errorY.FrameFilter {
	target := []byte("gin.(*Context).Next")
	return func(frame []byte) bool {
		return bytes.Contains(frame, target)
	}
}
