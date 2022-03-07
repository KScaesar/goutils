package xTest

import (
	"fmt"
	"time"

	"github.com/KScaesar/goutils"
)

func FakeTimeNow(fakeTime string) func() time.Time {
	return func() time.Time {
		t, err := goutils.TimeParse(fakeTime, false)
		if err != nil {
			panic(fmt.Sprintf("fake time now: %v", err))
		}
		return t
	}
}
