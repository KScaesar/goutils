package identity

import (
	"github.com/Min-Feng/goutils/base"
)

type AccountParam struct {
	UserName string
	Password HashedPassword
	Email    string
}

func NewAccount(param *AccountParam) Account {
	return Account{
		UserID:   base.NewID(),
		UserName: param.UserName,
		Password: param.Password,
		Email:    param.Email,
	}
}

type Account struct {
	UserID   base.ID
	UserName string
	Password HashedPassword
	Email    string
}
