package goutils

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSortKind_UnmarshalText_failed_because_not_match_sort_kind(t *testing.T) {
	type payLoad struct {
		Sort struct {
			CreateAt SortKind `json:"create_at"`
		} `bson:"sort"`
	}

	jsonBody := []byte(`
{
  "sort": {
    "create_at": "ace"
  }
}`)

	p := new(payLoad)
	err := json.Unmarshal(jsonBody, p)
	assert.Error(t, err)
}

func TestSortKind_UnmarshalText_successful(t *testing.T) {
	type payLoad struct {
		Sort struct {
			CreateAt SortKind `json:"create_at"`
			UpdateAt SortKind `json:"update_at"`
			DeleteAt SortKind `json:"delete_at"`
		} `bson:"sort"`
	}

	jsonBody := []byte(`
{
  "sort": {
    "create_at": "Asc",
    "update_at": "deSc",
    "delete_at": "aSc"
  }
}`)

	p := new(payLoad)
	err := json.Unmarshal(jsonBody, p)
	assert.NoError(t, err)

	assert.Equal(t, SortAsc, p.Sort.CreateAt, "create_at")
	assert.Equal(t, SortDesc, p.Sort.UpdateAt, "update_at")
	assert.Equal(t, SortAsc, p.Sort.DeleteAt, "delete_at")
}
