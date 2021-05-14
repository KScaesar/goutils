package identity

import (
	"database/sql/driver"
)

func NewUUID() UUID {
	return UUID{}
}

type UUID struct {
}

func (u *UUID) String() string {
	return ""
}

func (u *UUID) Value() (driver.Value, error) {
	return nil, nil
}

func (u *UUID) Scan(v interface{}) error {
	return nil
}
