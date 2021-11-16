package identity

import (
	"encoding/json"
	"strings"
)

type AuthServiceFactory interface {
	CreateAuthService(Visitor) (AuthService, error)
}

type AuthService interface {
	VerifyPermission(...Permission) bool
	VerifyDataOwnership(...IsMatchDataOwnership) bool
}

func NewAlwaysSuccessfulAuthService() AuthService {
	return alwaysSuccessfulAuthService{}
}

type alwaysSuccessfulAuthService struct{}

func (alwaysSuccessfulAuthService) VerifyPermission(...Permission) bool { return true }

func (alwaysSuccessfulAuthService) VerifyDataOwnership(...IsMatchDataOwnership) bool { return true }

func NewPermission(action, data string) Permission {
	return Permission{action: action, data: data}
}

type Permission struct {
	action string
	data   string
}

func (p *Permission) UnmarshalJSON(bytes []byte) error {
	var s string
	err := json.Unmarshal(bytes, &s)
	if err != nil {
		return err
	}

	split := strings.Split(s, ":")
	p.action = split[0]
	p.data = split[1]
	return nil
}

func (p Permission) MarshalJSON() ([]byte, error) {
	s := p.String()
	return json.Marshal(s)
}

func (p Permission) String() string {
	return p.action + ":" + p.data
}

type Visitor interface {
	UserID() string
}

type visitor struct {
	userID string
}

func (user *visitor) UserID() string {
	return user.userID
}

func verifyDataOwnership(visitor Visitor, fns ...IsMatchDataOwnership) bool {
	for _, isMatch := range fns {
		if !isMatch(visitor) {
			return false
		}
	}
	return true
}

type IsMatchDataOwnership func(Visitor) bool

func DataOwnershipBelongUser(userID string) IsMatchDataOwnership {
	return func(user Visitor) bool {
		return user.UserID() == userID
	}
}
