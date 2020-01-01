package api

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/koltyakov/gosip"
)

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

// Conf receives custom request config definition, e.g. custom headers, custom OData mod
func (users *Users) Conf(config *RequestConfig) *Users {
	users.config = config
	return users
}

// Select adds $select OData modifier
func (users *Users) Select(oDataSelect string) *Users {
	users.modifiers.AddSelect(oDataSelect)
	return users
}

// Expand adds $expand OData modifier
func (users *Users) Expand(oDataExpand string) *Users {
	users.modifiers.AddExpand(oDataExpand)
	return users
}

// Filter adds $filter OData modifier
func (users *Users) Filter(oDataFilter string) *Users {
	users.modifiers.AddFilter(oDataFilter)
	return users
}

// Top adds $top OData modifier
func (users *Users) Top(oDataTop int) *Users {
	users.modifiers.AddTop(oDataTop)
	return users
}

// OrderBy adds $orderby OData modifier
func (users *Users) OrderBy(oDataOrderBy string, ascending bool) *Users {
	users.modifiers.AddOrderBy(oDataOrderBy, ascending)
	return users
}

// Get ...
func (users *Users) Get() (UsersResp, error) {
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
		fmt.Sprintf("%s('%s')", users.endpoint, url.QueryEscape(loginName)),
		users.config,
	)
}

// GetByEmail ...
func (users *Users) GetByEmail(email string) *User {
	return NewUser(
		users.client,
		fmt.Sprintf("%s/GetByEmail('%s')", users.endpoint, url.QueryEscape(email)),
		users.config,
	)
}

/* Response helpers */

// Data : to get typed data
func (usersResp *UsersResp) Data() []UsersResp {
	collection, _ := parseODataCollection(*usersResp)
	users := []UsersResp{}
	for _, user := range collection {
		users = append(users, UsersResp(user))
	}
	return users
}

// Unmarshal : to unmarshal to custom object
func (usersResp *UsersResp) Unmarshal(obj interface{}) error {
	// collection := parseODataCollection(*usersResp)
	// data, _ := json.Marshal(collection)
	data, _ := parseODataCollectionPlain(*usersResp)
	return json.Unmarshal(data, obj)
}
