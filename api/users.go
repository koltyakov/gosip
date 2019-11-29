package api

import (
	"fmt"
	"net/url"

	"github.com/koltyakov/gosip"
)

// Users ...
type Users struct {
	client    *gosip.SPClient
	config    *RequestConfig
	endpoint  string
	modifiers map[string]string
}

// NewUsers ...
func NewUsers(client *gosip.SPClient, endpoint string, config *RequestConfig) *Users {
	return &Users{
		client:   client,
		endpoint: endpoint,
		config:   config,
	}
}

// ToURL ...
func (users *Users) ToURL() string {
	apiURL, _ := url.Parse(users.endpoint)
	query := url.Values{}
	for k, v := range users.modifiers {
		query.Add(k, trimMultiline(v))
	}
	apiURL.RawQuery = query.Encode()
	return apiURL.String()
}

// Conf ...
func (users *Users) Conf(config *RequestConfig) *Users {
	users.config = config
	return users
}

// Select ...
func (users *Users) Select(oDataSelect string) *Users {
	if users.modifiers == nil {
		users.modifiers = make(map[string]string)
	}
	users.modifiers["$select"] = oDataSelect
	return users
}

// Expand ...
func (users *Users) Expand(oDataExpand string) *Users {
	if users.modifiers == nil {
		users.modifiers = make(map[string]string)
	}
	users.modifiers["$expand"] = oDataExpand
	return users
}

// Filter ...
func (users *Users) Filter(oDataFilter string) *Users {
	if users.modifiers == nil {
		users.modifiers = make(map[string]string)
	}
	users.modifiers["$filter"] = oDataFilter
	return users
}

// Top ...
func (users *Users) Top(oDataTop int) *Users {
	if users.modifiers == nil {
		users.modifiers = make(map[string]string)
	}
	users.modifiers["$top"] = fmt.Sprintf("%d", oDataTop)
	return users
}

// OrderBy ...
func (users *Users) OrderBy(oDataOrderBy string, ascending bool) *Users {
	direction := "asc"
	if !ascending {
		direction = "desc"
	}
	if users.modifiers == nil {
		users.modifiers = make(map[string]string)
	}
	if users.modifiers["$orderby"] != "" {
		users.modifiers["$orderby"] += ","
	}
	users.modifiers["$orderby"] += fmt.Sprintf("%s %s", oDataOrderBy, direction)
	return users
}

// Get ...
func (users *Users) Get() ([]byte, error) {
	sp := NewHTTPClient(users.client)
	return sp.Get(users.ToURL(), getConfHeaders(users.config))
}

// GetByID ...
func (users *Users) GetByID(userID int) *User {
	return NewUser(
		users.client,
		fmt.Sprintf("%s/GetById(%d)", users.endpoint, userID),
		users.config,
	)
}

// GetByLoginName ...
func (users *Users) GetByLoginName(loginName string) *User {
	return NewUser(
		users.client,
		fmt.Sprintf("%s('%s')", users.endpoint, loginName),
		users.config,
	)
}

// GetByEmail ...
func (users *Users) GetByEmail(email string) *User {
	return NewUser(
		users.client,
		fmt.Sprintf("%s/GetByEmail('%s')", users.endpoint, email),
		users.config,
	)
}
