package api

import (
	"bytes"
	"context"
	"fmt"

	"github.com/koltyakov/gosip"
)

//go:generate ggen -ent User -conf -mods Select,Expand -helpers Data,Normalized

// User represents SharePoint Site User API queryable object struct
// Always use NewUser constructor instead of &User{}
type User struct {
	client    *gosip.SPClient
	config    *RequestConfig
	endpoint  string
	modifiers *ODataMods
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
		client:    client,
		endpoint:  endpoint,
		config:    config,
		modifiers: NewODataMods(),
	}
}

// ToURL gets endpoint with modificators raw URL
func (user *User) ToURL() string {
	return toURL(user.endpoint, user.modifiers)
}

// Get gets this user data object
func (user *User) Get(ctx context.Context) (UserResp, error) {
	client := NewHTTPClient(user.client)
	return client.Get(ctx, user.ToURL(), user.config)
}

// Update updates User's metadata with properties provided in `body` parameter
// where `body` is byte array representation of JSON string payload relevant to SP.User object
func (user *User) Update(ctx context.Context, body []byte) (UserResp, error) {
	body = patchMetadataType(body, "SP.User")
	client := NewHTTPClient(user.client)
	return client.Update(ctx, user.endpoint, bytes.NewBuffer(body), user.config)
}

// Groups gets Groups API instance queryable collection for this User
func (user *User) Groups() *Groups {
	return NewGroups(
		user.client,
		fmt.Sprintf("%s/Groups", user.endpoint),
		user.config,
	)
}
