package api

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/koltyakov/gosip"
)

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

// FromURL gets User object using its API URL
func (user *User) FromURL(url string) *User {
	url = strings.Split(url, "?")[0]
	return NewUser(user.client, url, user.config)
}

// Conf receives custom request config definition, e.g. custom headers, custom OData mod
func (user *User) Conf(config *RequestConfig) *User {
	user.config = config
	return user
}

// Select adds $select OData modifier
func (user *User) Select(oDataSelect string) *User {
	user.modifiers.AddSelect(oDataSelect)
	return user
}

// Expand adds $expand OData modifier
func (user *User) Expand(oDataExpand string) *User {
	user.modifiers.AddExpand(oDataExpand)
	return user
}

// Get gets this user data object
func (user *User) Get() (UserResp, error) {
	sp := NewHTTPClient(user.client)
	return sp.Get(user.ToURL(), getConfHeaders(user.config))
}

// Update updates User's metadata with properties provided in `body` parameter
// where `body` is byte array representation of JSON string payload relevalt to SP.User object
func (user *User) Update(body []byte) (UserResp, error) {
	body = patchMetadataType(body, "SP.User")
	sp := NewHTTPClient(user.client)
	return sp.Update(user.endpoint, body, getConfHeaders(user.config))
}

// Groups gets Groups API instance queryable collection for this User
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
	data := NormalizeODataItem(*userResp)
	res := &UserInfo{}
	json.Unmarshal(data, &res)
	return res
}

// Normalized returns normalized body
func (userResp *UserResp) Normalized() []byte {
	return NormalizeODataItem(*userResp)
}
