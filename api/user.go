package api

import (
	"fmt"
	"net/url"

	"github.com/koltyakov/gosip"
)

// User ...
type User struct {
	client    *gosip.SPClient
	config    *RequestConfig
	endpoint  string
	modifiers map[string]string
}

// UserInfo ...
type UserInfo struct {
	Email         string `json:"Email"`
	ID            int    `json:"Id"`
	IsHiddenInUI  bool   `json:"IsHiddenInUI"`
	IsSiteAdmin   bool   `json:"IsSiteAdmin"`
	LoginName     string `json:"LoginName"`
	PrincipalType int    `json:"PrincipalType"`
	Title         int    `json:"Title"`
}

// ToURL ...
func (user *User) ToURL() string {
	return user.endpoint
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
func (user *User) Get() ([]byte, error) {
	apiURL, _ := url.Parse(user.endpoint)
	query := url.Values{}
	for k, v := range user.modifiers {
		query.Add(k, trimMultiline(v))
	}
	apiURL.RawQuery = query.Encode()
	sp := NewHTTPClient(user.client)
	return sp.Get(apiURL.String(), getConfHeaders(user.config))
}

// Groups ...
func (user *User) Groups() *Groups {
	return &Groups{
		client: user.client,
		config: user.config,
		endpoint: fmt.Sprintf(
			"%s/Groups",
			user.endpoint,
		),
	}
}
