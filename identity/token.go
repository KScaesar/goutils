package identity

import (
	"time"

	"github.com/Min-Feng/goutils"
)

var (
	TokenKindAccess  = "Access"
	TokenKindRefresh = "Refresh"
)

type TokenConfig struct {
	AccessInterval  time.Duration
	RefreshInterval time.Duration
}

func NewAccessToken(id goutils.ID) Token {
	return Token{
		Kind:     TokenKindAccess,
		UserID:   id,
		DeadLine: TimeNow().Add(setting.Token.AccessInterval),
	}
}

func NewRefreshToken(id goutils.ID) Token {
	return Token{
		Kind:     TokenKindRefresh,
		UserID:   id,
		DeadLine: TimeNow().Add(setting.Token.RefreshInterval),
	}
}

type Token struct {
	Kind     string
	UserID   goutils.ID
	DeadLine time.Time
}

func (t *Token) ResetDeadLine() {
	switch t.Kind {
	case TokenKindAccess:
		t.DeadLine = TimeNow().Add(setting.Token.AccessInterval)
	case TokenKindRefresh:
		t.DeadLine = TimeNow().Add(setting.Token.RefreshInterval)
	}
}

func (t Token) IsExpired() bool {
	return TimeNow().After(t.DeadLine)
}
