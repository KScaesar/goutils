package goutils

import (
	"bytes"
	"fmt"
)

const (
	SortNone SortKind = ""
	SortDesc SortKind = "desc"
	SortAsc  SortKind = "asc"
)

type SortKind string

func (s *SortKind) UnmarshalText(text []byte) error {
	value := SortKind(bytes.ToLower(text))

	switch value {
	case SortNone, SortDesc, SortAsc:
		*s = value
		return nil
	default:
		return fmt.Errorf("not match sort kind: value = %v", string(text))
	}
}
