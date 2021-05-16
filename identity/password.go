package identity

import (
	"database/sql"
	"database/sql/driver"
	"encoding"
	"fmt"
	"unsafe"

	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"

	"github.com/Min-Feng/goutils/errorY"
)

var (
	_ driver.Valuer            = (*HashedPassword)(nil)
	_ sql.Scanner              = (*HashedPassword)(nil)
	_ encoding.TextMarshaler   = (*HashedPassword)(nil)
	_ encoding.TextUnmarshaler = (*HashedPassword)(nil)
	_ fmt.Stringer             = (*HashedPassword)(nil)
)

type HashedPassword struct {
	bytes []byte
}

func (pw HashedPassword) VerifyPassword(plain PlainPassword) error {
	err := bcrypt.CompareHashAndPassword(pw.bytes, plain.bytes)
	if err != nil {
		return errorY.Wrap(ErrAuthentication, err.Error())
	}
	return nil
}

// Value 實現 driver.Value,
// 但疑問 回傳 []byte, string 的差異?
func (pw HashedPassword) Value() (driver.Value, error) {
	return pw.bytes, nil
}

func (pw *HashedPassword) Scan(src interface{}) error {
	switch v := src.(type) {
	case []byte:
		pw.bytes = v
	case string:
		pw.bytes = []byte(v)
	}
	return nil
}

func (pw HashedPassword) MarshalText() (text []byte, err error) {
	return pw.bytes, nil
}

func (pw *HashedPassword) UnmarshalText(text []byte) error {
	pw.bytes = text
	return nil
}

func (pw HashedPassword) String() string {
	return *(*string)(unsafe.Pointer(&pw.bytes))
}

func NewPlainPassword(plainPW string) PlainPassword {
	return PlainPassword{
		// 因為這個 []byte 會被送到其他函數進行操作, 所以不能進行 string to []byte 的特例優化
		// bytes:  *(*[]byte)(unsafe.Pointer(&plainPW)),
		bytes:  []byte(plainPW),
		string: plainPW,
	}
}

type PlainPassword struct {
	bytes  []byte
	string string
}

func (pw PlainPassword) Bcrypt() (HashedPassword, error) {
	if err := pw.rule(); err != nil {
		return HashedPassword{}, errorY.WrapMessage(err, "violation of rules")
	}

	const cost = 10
	hash, err := bcrypt.GenerateFromPassword(pw.bytes, cost)
	if err != nil {
		return HashedPassword{}, errorY.Wrap(errorY.ErrSystem, err.Error())
	}

	return HashedPassword{
		bytes: hash,
	}, nil
}

func (pw *PlainPassword) UnmarshalText(text []byte) error {
	pw.bytes = text
	pw.string = *(*string)(unsafe.Pointer(&text))
	return pw.rule()
}

func (pw PlainPassword) String() string {
	return pw.string
}

func (pw PlainPassword) rule() error {
	var plainPWRuler = validator.New()
	err := plainPWRuler.Var(pw.string, "gte=8,alphanum")
	if err != nil {
		return errorY.Wrap(errorY.ErrInvalidParams, err.Error())
	}
	return nil
}
