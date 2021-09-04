package goutils

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTimeParse(t *testing.T) {

	tests := []struct {
		name      string
		timeValue string
	}{
		{timeValue: "2020-09-04 20:30:30 +0800"},
		{timeValue: "2020-09-04 20:30:30 -05:30"},
		{timeValue: "2020-09-04 20:30:30 +04:00"},
		// {timeValue: "2020-09-04 20:30:30Z"}, // 搞不懂 Z 時區的表示法
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			value, err := TimeParse(tt.timeValue)
			assert.NoError(t, err)

			fmt.Println(value.UTC().Format("2006-01-02 15:04:05 -07:00"))
		})
	}
}
