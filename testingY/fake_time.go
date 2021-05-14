package testingY

import (
	"fmt"
	"time"

	"github.com/Min-Feng/goutils/kit"
)

func FakeTimeNow(fakeTime string) func() time.Time {
	return func() time.Time {
		t, err := kit.TimeParse(fakeTime)
		if err != nil {
			panic(fmt.Sprintf("fake time now: %v", err))
		}
		return t
	}
}
