package identity

import (
	"github.com/Min-Feng/goutils"
)

type AccountParam struct {
	UserName string
	Password HashedPassword
	Email    string
}

func NewAccount(param *AccountParam) Account {
	return Account{
		UserID:   goutils.NewID(),
		UserName: param.UserName,
		Password: param.Password,
		Email:    param.Email,
	}
}

type Account struct {
	UserID   goutils.ID
	UserName string
	Password HashedPassword
	Email    string
}
