package api

import (
	"fmt"
	"net/url"

	"github.com/koltyakov/gosip"
)

//go:generate ggen -ent Users -conf -mods Select,Expand,Filter,Top,OrderBy

// Users represent SharePoint Site Users API queryable collection struct
// Always use NewUsers constructor instead of &Users{}
type Users struct {
	client    *gosip.SPClient
	config    *RequestConfig
	endpoint  string
	modifiers *ODataMods
}

// UsersResp - site users response type with helper processor methods
type UsersResp []byte

// NewUsers - Users struct constructor function
func NewUsers(client *gosip.SPClient, endpoint string, config *RequestConfig) *Users {
	return &Users{
		client:    client,
		endpoint:  endpoint,
		config:    config,
		modifiers: NewODataMods(),
	}
}

// ToURL gets endpoint with modificators raw URL
func (users *Users) ToURL() string {
	return toURL(users.endpoint, users.modifiers)
}

// Get gets Users queryable collection
func (users *Users) Get() (UsersResp, error) {
	sp := NewHTTPClient(users.client)
	return sp.Get(users.ToURL(), getConfHeaders(users.config))
}

// GetByID gets a user by his/her ID (numeric ID from User Information List)
func (users *Users) GetByID(userID int) *User {
	return NewUser(
		users.client,
		fmt.Sprintf("%s/GetById(%d)", users.endpoint, userID),
		users.config,
	)
}

// GetByLoginName gets a user by his/her login name
func (users *Users) GetByLoginName(loginName string) *User {
	return NewUser(
		users.client,
		fmt.Sprintf("%s('%s')", users.endpoint, url.QueryEscape(loginName)),
		users.config,
	)
}

// GetByEmail gets a user by his/her email address
func (users *Users) GetByEmail(email string) *User {
	return NewUser(
		users.client,
		fmt.Sprintf("%s/GetByEmail('%s')", users.endpoint, url.QueryEscape(email)),
		users.config,
	)
}

/* Response helpers */

// Data : to get typed data
func (usersResp *UsersResp) Data() []UserResp {
	collection, _ := normalizeODataCollection(*usersResp)
	users := []UserResp{}
	for _, user := range collection {
		users = append(users, UserResp(user))
	}
	return users
}

// Normalized returns normalized body
func (usersResp *UsersResp) Normalized() []byte {
	normalized, _ := NormalizeODataCollection(*usersResp)
	return normalized
}
