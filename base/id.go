package base

import (
	"database/sql/driver"
)

func NewID() ID {
	return ID{}
}

type ID struct {
}

func (u *ID) String() string {
	return ""
}

func (u *ID) Value() (driver.Value, error) {
	return nil, nil
}

func (u *ID) Scan(v interface{}) error {
	return nil
}
