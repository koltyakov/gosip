package api

import (
	"fmt"
	"net/url"

	"github.com/koltyakov/gosip"
)

// Users ...
type Users struct {
	client   *gosip.SPClient
	conf     *Conf
	endpoint string
	oSelect  string
	oExpand  string
	oFilter  string
	oTop     int
	oOrderBy string
}

// Conf ...
func (users *Users) Conf(conf *Conf) *Users {
	users.conf = conf
	return users
}

// Select ...
func (users *Users) Select(oDataSelect string) *Users {
	users.oSelect = oDataSelect
	return users
}

// Expand ...
func (users *Users) Expand(oDataExpand string) *Users {
	users.oExpand = oDataExpand
	return users
}

// Filter ...
func (users *Users) Filter(oDataFilter string) *Users {
	users.oFilter = oDataFilter
	return users
}

// Top ...
func (users *Users) Top(oDataTop int) *Users {
	users.oTop = oDataTop
	return users
}

// OrderBy ...
func (users *Users) OrderBy(oDataOrderBy string, ascending bool) *Users {
	direction := "asc"
	if !ascending {
		direction = "desc"
	}
	if users.oOrderBy != "" {
		users.oOrderBy += ","
	}
	users.oOrderBy += fmt.Sprintf("%s %s", oDataOrderBy, direction)
	return users
}

// Get ...
func (users *Users) Get() ([]byte, error) {
	apiURL, _ := url.Parse(users.endpoint)
	query := url.Values{}
	if users.oSelect != "" {
		query.Add("$select", TrimMultiline(users.oSelect))
	}
	if users.oExpand != "" {
		query.Add("$expand", TrimMultiline(users.oExpand))
	}
	if users.oFilter != "" {
		query.Add("$filter", TrimMultiline(users.oFilter))
	}
	if users.oTop != 0 {
		query.Add("$top", fmt.Sprintf("%d", users.oTop))
	}
	if users.oOrderBy != "" {
		query.Add("$orderBy", TrimMultiline(users.oOrderBy))
	}
	apiURL.RawQuery = query.Encode()
	sp := &HTTPClient{SPClient: users.client}
	return sp.Get(apiURL.String(), GetConfHeaders(users.conf))
}

// GetByID ...
func (users *Users) GetByID(userID int) *User {
	return &User{
		client: users.client,
		conf:   users.conf,
		endpoint: fmt.Sprintf("%s/GetById(%d)",
			users.endpoint,
			userID,
		),
	}
}

// GetByName ...
func (users *Users) GetByLoginName(loginName string) *User {
	return &User{
		client: users.client,
		conf:   users.conf,
		endpoint: fmt.Sprintf("%s('%s')",
			users.endpoint,
			loginName,
		),
	}
}

// GetByName ...
func (users *Users) GetByEmail(email string) *User {
	return &User{
		client: users.client,
		conf:   users.conf,
		endpoint: fmt.Sprintf("%s/GetByEmail('%s')",
			users.endpoint,
			email,
		),
	}
}
