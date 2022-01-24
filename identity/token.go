package identity

import (
	"time"
)

type TokenKind string

const (
	TokenKindAccess  TokenKind = "Access"
	TokenKindRefresh TokenKind = "Refresh"
)

type TokenConfig struct {
	AccessInterval  time.Duration
	RefreshInterval time.Duration
}

func NewAccessToken(userID string) Token {
	return Token{
		Kind:     TokenKindAccess,
		UserID:   userID,
		DeadLine: TimeNow().Add(setting.Token.AccessInterval),
	}
}

func NewRefreshToken(userID string) Token {
	return Token{
		Kind:     TokenKindRefresh,
		UserID:   userID,
		DeadLine: TimeNow().Add(setting.Token.RefreshInterval),
	}
}

type Token struct {
	Kind     TokenKind
	UserID   string
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
