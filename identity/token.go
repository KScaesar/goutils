package identity

import (
	"time"

	"github.com/Min-Feng/goutils/base"
)

type TokenConfig struct {
	AccessInterval  time.Duration
	RefreshInterval time.Duration
}

func NewAccessToken(id base.ID) Token {
	return Token{
		UserID:   id,
		DeadLine: TimeNow().Add(setting.Token.AccessInterval),
	}
}

func NewRefreshToken(id base.ID) Token {
	return Token{
		UserID:   id,
		DeadLine: TimeNow().Add(setting.Token.RefreshInterval),
	}
}

type Token struct {
	UserID   base.ID
	DeadLine time.Time
}
