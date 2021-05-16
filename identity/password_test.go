package identity

import (
	"encoding/json"
	"fmt"
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

func TestHashedPassword_MarshalText(t *testing.T) {
	hash, err := NewPlainPassword("qazwsx123").Bcrypt()
	assert.NoError(t, err)

	type PayLoad struct {
		HashedPW HashedPassword `json:"password"`
	}
	actualJson, err := json.Marshal(PayLoad{HashedPW: hash})
	assert.NoError(t, err)

	expected := []byte(fmt.Sprintf(`{"password":"%v"}`, hash.String()))
	assert.Equal(t, expected, actualJson)
}

func TestHashedPassword_Verify(t *testing.T) {
	plainPW := "qazwsx123"
	hash, err := NewPlainPassword(plainPW).Bcrypt()
	assert.NoError(t, err)

	plain := NewPlainPassword(plainPW)
	err = hash.VerifyPassword(plain)
	assert.NoError(t, err)
}
