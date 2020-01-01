package api

import (
	"encoding/json"
	"fmt"

	"github.com/koltyakov/gosip"
)

// User represents SharePoint Site User API queryable object struct
type User struct {
	client    *gosip.SPClient
	config    *RequestConfig
	endpoint  string
	modifiers map[string]string
}

// UserInfo - site user API response payload structure
type UserInfo struct {
	Email         string `json:"Email"`
	ID            int    `json:"Id"`
	IsHiddenInUI  bool   `json:"IsHiddenInUI"`
	IsSiteAdmin   bool   `json:"IsSiteAdmin"`
	LoginName     string `json:"LoginName"`
	PrincipalType int    `json:"PrincipalType"`
	Title         string `json:"Title"`
}

// UserResp - site user response type with helper processor methods
type UserResp []byte

// NewUser - User struct constructor function
func NewUser(client *gosip.SPClient, endpoint string, config *RequestConfig) *User {
	return &User{
		client:   client,
		endpoint: endpoint,
		config:   config,
	}
}

// ToURL gets endpoint with modificators raw URL ...
func (user *User) ToURL() string {
	return toURL(user.endpoint, user.modifiers)
}

// Conf ...
func (user *User) Conf(config *RequestConfig) *User {
	user.config = config
	return user
}

// Select ...
func (user *User) Select(oDataSelect string) *User {
	if user.modifiers == nil {
		user.modifiers = make(map[string]string)
	}
	user.modifiers["$select"] = oDataSelect
	return user
}

// Expand ...
func (user *User) Expand(oDataExpand string) *User {
	if user.modifiers == nil {
		user.modifiers = make(map[string]string)
	}
	user.modifiers["$expand"] = oDataExpand
	return user
}

// Get ...
func (user *User) Get() (UserResp, error) {
	sp := NewHTTPClient(user.client)
	return sp.Get(user.ToURL(), getConfHeaders(user.config))
}

// Groups ...
func (user *User) Groups() *Groups {
	return NewGroups(
		user.client,
		fmt.Sprintf("%s/Groups", user.endpoint),
		user.config,
	)
}

/* Response helpers */

// Data : to get typed data
func (userResp *UserResp) Data() *UserInfo {
	data := parseODataItem(*userResp)
	res := &UserInfo{}
	json.Unmarshal(data, &res)
	return res
}

// Unmarshal : to unmarshal to custom object
func (userResp *UserResp) Unmarshal(obj interface{}) error {
	data := parseODataItem(*userResp)
	return json.Unmarshal(data, obj)
}
