package identity

import (
	"context"
	"encoding/json"
	"strings"
)

type CheckServiceFactory interface {
	CreateCheckService(ctx context.Context, enable bool) (CheckService, error)
}

type CheckService interface {
	PermissionOk(...Permission) bool
	OwnershipOk(...DataAttribute) bool
}

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
	p.data = split[0]
	p.action = split[1]
	return nil
}

func (p Permission) MarshalJSON() ([]byte, error) {
	s := p.String()
	return json.Marshal(s)
}

func (p Permission) String() string {
	return p.action + ":" + p.data
}

func NewOwnershipSet(loginUserData ...DataAttribute) OwnershipSet {
	set := make(OwnershipSet, len(loginUserData))
	for _, attribute := range loginUserData {
		set[attribute.Name] = attribute.Value
	}
	return set
}

type OwnershipSet map[AttributeKey]interface{}

func (set OwnershipSet) OwnershipOk(attributes ...DataAttribute) bool {
	for _, atr := range attributes {
		if !set.matchAttribute(atr) {
			return false
		}
	}
	return true
}

func (set OwnershipSet) matchAttribute(atr DataAttribute) bool {
	return set[atr.Name] == atr.Value
}

type AttributeKey string

type DataAttribute struct {
	Name  AttributeKey
	Value interface{}
}

const (
	AttributeKeyUser AttributeKey = "user_id"
)

func AttributeUserID(userID string) DataAttribute {
	return DataAttribute{
		Name:  AttributeKeyUser,
		Value: userID,
	}
}
