package identity

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPlainPassword_rule_Success(t *testing.T) {
	buf := []byte(`{"password":"qazwsx123"}`)
	type PayLoad struct {
		Password PlainPassword `json:"password"`
	}
	payLoad := new(PayLoad)
	err := json.Unmarshal(buf, payLoad)
	assert.NoError(t, err)
}

func TestPlainPassword_rule_Failed_Length_Less_8(t *testing.T) {
	buf := []byte(`{"password":"qaZ23"}`)
	type PayLoad struct {
		Password PlainPassword `json:"password"`
	}
	payLoad := new(PayLoad)
	err := json.Unmarshal(buf, payLoad)
	assert.Error(t, err)
}

func TestPlainPassword_Bcrypt(t *testing.T) {
	hash1, _ := NewPlainPassword("123QWEasd").Bcrypt()
	hash2, _ := NewPlainPassword("123QWEasd").Bcrypt()
	assert.NotEqualf(t, hash2, hash1, "同樣的明碼, 每次產生的 hash 都不同")

	plain := "123QWEasd"
	assert.True(t, hash1.VerifyPassword(plain))
	assert.True(t, hash2.VerifyPassword(plain))
}

func TestHashedPassword_MarshalText(t *testing.T) {
	// 保存的值 是 $2a$10$OoYhOY9FHNTT4N1n5F2L3eqM9TkilUm5FKf0KF1RI53SqdRbXficu
	// 不包含 雙引號
	hash := HashedPassword{
		bytes: []byte("$2a$10$OoYhOY9FHNTT4N1n5F2L3eqM9TkilUm5FKf0KF1RI53SqdRbXficu"),
	}

	type PayLoad struct {
		HashedPW HashedPassword `json:"password"`
	}
	actualJson, err := json.Marshal(PayLoad{HashedPW: hash})
	assert.NoError(t, err)

	expected := []byte(`{"password":"$2a$10$OoYhOY9FHNTT4N1n5F2L3eqM9TkilUm5FKf0KF1RI53SqdRbXficu"}`)
	assert.Equal(t, expected, actualJson)
}

func TestHashedPassword_VerifyPassword(t *testing.T) {
	type PayLoad struct {
		Password string `json:"password"`
	}
	buf := []byte(`{"password":"qazwsx123"}`)
	payLoad := new(PayLoad)
	err := json.Unmarshal(buf, payLoad)
	assert.NoError(t, err)

	hash := HashedPassword{
		bytes: []byte("$2a$10$8wfdlrkWi2QnZtfIs6jIWOZjGW3r6SVzTMv0O83JrGG42xYwddLym"),
	}
	assert.True(t, hash.VerifyPassword(payLoad.Password))
}
