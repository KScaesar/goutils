package goutils

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewPageResponse_json_Marshal(t *testing.T) {
	opt := PageOption{
		PageNumber: 4,
		PageSize:   9,
	}

	expected := `{"totalPage":4,"totalCount":34,"targetPageNumber":4,"targetPageSize":7}`
	body, _ := json.Marshal(NewPageResponse(opt, 34))
	assert.Equal(t, expected, string(body))
}

func TestNewPageResponse(t *testing.T) {
	type args struct {
		opt        PageOption
		totalCount int64
	}
	tests := []struct {
		name   string
		args   args
		expect PageResponse
	}{
		{
			name: "1",
			args: args{
				opt:        NewPageOption(1, 10),
				totalCount: 34,
			},
			expect: PageResponse{
				TotalPageNumber:  4,
				TotalCount:       34,
				TargetPageNumber: 1,
				TargetPageSize:   10,
			},
		},
		{
			name: "2",
			args: args{
				opt:        NewPageOption(2, 40),
				totalCount: 34,
			},
			expect: PageResponse{
				TotalPageNumber:  1,
				TotalCount:       34,
				TargetPageNumber: 2,
				TargetPageSize:   0,
			},
		},
		{
			name: "3",
			args: args{
				opt:        NewPageOption(2, 30),
				totalCount: 34,
			},
			expect: PageResponse{
				TotalPageNumber:  2,
				TotalCount:       34,
				TargetPageNumber: 2,
				TargetPageSize:   4,
			},
		},
		{
			name: "4",
			args: args{
				opt:        NewPageOption(1, 50),
				totalCount: 34,
			},
			expect: PageResponse{
				TotalPageNumber:  1,
				TotalCount:       34,
				TargetPageNumber: 1,
				TargetPageSize:   34,
			},
		},
		{
			name: "5",
			args: args{
				opt:        NewPageOption(1, 40),
				totalCount: 40,
			},
			expect: PageResponse{
				TotalPageNumber:  1,
				TotalCount:       40,
				TargetPageNumber: 1,
				TargetPageSize:   40,
			},
		},
		{
			name: "6",
			args: args{
				opt:        NewPageOption(1, 40),
				totalCount: 0,
			},
			expect: PageResponse{
				TotalPageNumber:  1,
				TotalCount:       0,
				TargetPageNumber: 1,
				TargetPageSize:   0,
			},
		},
		{
			name: "7",
			args: args{
				opt:        NewPageOption(0, 0),
				totalCount: 0,
			},
			expect: PageResponse{
				TotalPageNumber:  1,
				TotalCount:       0,
				TargetPageNumber: 1,
				TargetPageSize:   0,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			actual := NewPageResponse(tt.args.opt, tt.args.totalCount)
			assert.Equal(t, tt.expect, *actual)
		})
	}
}
