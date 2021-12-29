package goutils

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Min-Feng/goutils/xLog"
)

func TestTimeParse(t *testing.T) {
	// xLog.SetGlobalLevel("panic")

	tests := []struct {
		name      string
		timeValue string
	}{
		{timeValue: "2020-09-04 20:30:30-05:30"},
		{timeValue: "2020-09-04 20:30:30 -05:30"},
		{timeValue: "2020-09-04 20:30:30+04:00"},
		{timeValue: "2020-09-04 20:30:30 +04:00"},
		{timeValue: "2020-09-04 20:30:30"}, // 沒加上時區, 都視為 UTC 時區
		{timeValue: "2020-09-04 20:30:30Z"},
		{timeValue: "2020-10-17"}, // 沒加上時區, 都視為 UTC 時區

		// Z 的用途 和 T 一樣, 標示間隔 同時保有字串連續性
		// 且同等正負號的位置
		// {timeValue: "2020-09-04 20:30:30Z"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			value, err := TimeParse(tt.timeValue)
			assert.NoError(t, err)

			xLog.Debug().
				Str("input", tt.timeValue).
				Time("log_format", value).
				Interface("my_format", Time(value)).
				Send()
		})
	}
}
