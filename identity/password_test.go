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

func TestHashedPassword_MarshalText(t *testing.T) {
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

func TestHashedPassword_Verify(t *testing.T) {
	hash := HashedPassword{
		bytes: []byte("$2a$10$8wfdlrkWi2QnZtfIs6jIWOZjGW3r6SVzTMv0O83JrGG42xYwddLym"),
	}

	plainPW := "qazwsx123"
	plain := NewPlainPassword(plainPW)
	err := hash.VerifyPassword(plain)
	assert.NoError(t, err)
}
