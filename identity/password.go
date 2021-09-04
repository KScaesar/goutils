package identity

import (
	"database/sql"
	"database/sql/driver"
	"encoding"
	"fmt"
	"unsafe"

	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"

	"github.com/Min-Feng/goutils/errors"
)

var (
	_ driver.Valuer            = (*HashedPassword)(nil)
	_ sql.Scanner              = (*HashedPassword)(nil)
	_ encoding.TextMarshaler   = (*HashedPassword)(nil)
	_ encoding.TextUnmarshaler = (*HashedPassword)(nil)
	_ fmt.Stringer             = (*HashedPassword)(nil)
)

// HashedPassword is generated from PlainPassword.Bcrypt
type HashedPassword struct {
	bytes []byte
}

func (pw HashedPassword) VerifyPassword(plainPW string) bool {
	err := bcrypt.CompareHashAndPassword(pw.bytes, []byte(plainPW))
	if err != nil {
		return false
	}
	return true
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
	pw := PlainPassword{
		// 因為這個 []byte 會被送到其他函數進行操作
		// 所以不能進行 string to []byte 的特例優化, 因為字串的不可變性
		// bytes:  *(*[]byte)(unsafe.Pointer(&plainPW))
		bytes:  []byte(plainPW),
		string: plainPW,
	}

	// 猶豫是否應該回傳 error
	// 但會讓 api 不好用, 無法達成類似這樣的呼叫, NewPlainPassword("123QWEasd").Bcrypt()
	// 如果只有一個檢查點 CheckRule(), 感覺可以先不用回傳 error
	// 只要在少數幾個 method 手動呼叫 CheckRule()
	// 如果要檢查的東西太多了, 或許就應該好好回傳 error
	// pw.err = pw.CheckRule()

	return pw
}

type PlainPassword struct {
	bytes  []byte
	string string
}

func (pw PlainPassword) Bcrypt() (HashedPassword, error) {
	if err := pw.CheckRule(); err != nil {
		return HashedPassword{}, err
	}

	const cost = 10
	hash, err := bcrypt.GenerateFromPassword(pw.bytes, cost)
	if err != nil {
		return HashedPassword{}, errors.Wrap(errors.ErrSystem, err.Error())
	}

	return HashedPassword{
		bytes: hash,
	}, nil
}

func (pw PlainPassword) CheckRule() (Err error) {
	defer func() {
		if Err != nil {
			Err = errors.Wrap(errors.ErrInvalidParams, "violation of password rules: %v", Err)
		}
	}()

	plainPWRuler := validator.New()

	if err := plainPWRuler.Var(pw.string, "alphanum"); err != nil {
		return err
	}
	if err := plainPWRuler.Var(pw.string, "min=8,max=24"); err != nil {
		return err
	}
	return nil
}

func (pw *PlainPassword) UnmarshalText(text []byte) error {
	pw.bytes = text
	pw.string = *(*string)(unsafe.Pointer(&text))
	return pw.CheckRule()
}

func (pw PlainPassword) String() string {
	if err := pw.CheckRule(); err != nil {
		return fmt.Sprintf("%v: %v", pw.string, err)
	}
	return pw.string
}
