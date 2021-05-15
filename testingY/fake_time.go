package testingY

import (
	"fmt"
	"time"

	"github.com/Min-Feng/goutils/kits"
)

func FakeTimeNow(fakeTime string) func() time.Time {
	return func() time.Time {
		t, err := kits.TimeParse(fakeTime)
		if err != nil {
			panic(fmt.Sprintf("fake time now: %v", err))
		}
		return t
	}
}
