package base

import (
	"time"

	"github.com/Min-Feng/goutils/errorY"
)

var timeSpec = []string{
	"2006-01-02 15:04:05",
	"2006-01-02",
	"2006-01-02 15:04",
	"2006-01-02 15:04:05Z07:00",
	"2006-01-02T15:04",
	"2006-01-02T15:04:05",
	time.RFC3339,
}

func TimeParse(timeLayout string) (t time.Time, err error) {
	for _, spec := range timeSpec {
		t, err = time.ParseInLocation(spec, timeLayout, time.Local)
		if err == nil {
			return t, nil
		}
	}
	return time.Time{}, errorY.Wrap(errorY.ErrSystem, err.Error())
}
