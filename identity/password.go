package identity

import (
	"bytes"
	"database/sql"
	"database/sql/driver"
	"encoding"
	"fmt"

	"github.com/Min-Feng/goutils/errorY"
)

type PasswordConfig struct {
	Key  []byte
	Salt []byte
}

func EncryptPassword(cfg PasswordConfig, password []byte) (encrypt []byte, err error) {
	return nil, nil
}

var (
	_ driver.Valuer            = (*Password)(nil)
	_ sql.Scanner              = (*Password)(nil)
	_ encoding.TextMarshaler   = (*Password)(nil)
	_ encoding.TextUnmarshaler = (*Password)(nil)
	_ fmt.Stringer             = (*Password)(nil)
)

func NewPassword(text []byte) (Password, error) {
	encryptPassword, err := EncryptPassword(setting.Password, text)
	if err != nil {
		return Password{}, errorY.WrapMessage(err, "encrypt password")
	}

	return Password{
		encrypt: encryptPassword,
	}, nil
}

type Password struct {
	encrypt []byte
}

func (p Password) Equal(other Password) bool {
	return bytes.Equal(p.encrypt, other.encrypt)
}

func (p Password) Value() (driver.Value, error) {
	return p.encrypt, nil
}

func (p *Password) Scan(src interface{}) error {
	switch v := src.(type) {
	case []byte:
		p.encrypt = v
	case string:
		p.encrypt = []byte(v)
	}
	return nil
}

func (p Password) MarshalText() (text []byte, err error) {
	return p.encrypt, nil
}

func (p *Password) UnmarshalText(text []byte) error {
	encryptPassword, err := EncryptPassword(setting.Password, text)
	if err != nil {
		return errorY.WrapMessage(err, "encrypt password")
	}

	p.encrypt = encryptPassword
	return nil
}

func (p Password) String() string {
	return string(p.encrypt)
}
