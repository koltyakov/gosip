package api

import (
	"fmt"
	"net/url"

	"github.com/koltyakov/gosip"
)

// User ...
type User struct {
	client   *gosip.SPClient
	conf     *Conf
	endpoint string
	oSelect  string
	oExpand  string
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

// Conf ...
func (user *User) Conf(conf *Conf) *User {
	user.conf = conf
	return user
}

// Select ...
func (user *User) Select(oDataSelect string) *User {
	user.oSelect = oDataSelect
	return user
}

// Expand ...
func (user *User) Expand(oDataExpand string) *User {
	user.oExpand = oDataExpand
	return user
}

// Get ...
func (user *User) Get() ([]byte, error) {
	apiURL, _ := url.Parse(user.endpoint)
	query := url.Values{}
	if user.oSelect != "" {
		query.Add("$select", trimMultiline(user.oSelect))
	}
	if user.oExpand != "" {
		query.Add("$expand", trimMultiline(user.oExpand))
	}
	apiURL.RawQuery = query.Encode()
	sp := NewHTTPClient(user.client)
	return sp.Get(apiURL.String(), getConfHeaders(user.conf))
}

// Groups ...
func (user *User) Groups() *Groups {
	return &Groups{
		client: user.client,
		conf:   user.conf,
		endpoint: fmt.Sprintf(
			"%s/Groups",
			user.endpoint,
		),
	}
}
