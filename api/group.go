package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/koltyakov/gosip"
)

//go:generate ggen -ent Group -conf -mods Select,Expand -helpers Data,Normalized

// Group represents SharePoint Site Groups API queryable object struct
// Always use NewGroup constructor instead of &Group{}
type Group struct {
	client    *gosip.SPClient
	config    *RequestConfig
	endpoint  string
	modifiers *ODataMods
}

// GroupInfo - site group API response payload structure
type GroupInfo struct {
	AllowMembersEditMembership     bool   `json:"AllowMembersEditMembership"`
	AllowRequestToJoinLeave        bool   `json:"AllowRequestToJoinLeave"`
	AutoAcceptRequestToJoinLeave   bool   `json:"AutoAcceptRequestToJoinLeave"`
	Description                    string `json:"Description"`
	ID                             int    `json:"Id"`
	IsHiddenInUI                   bool   `json:"IsHiddenInUI"`
	LoginName                      string `json:"LoginName"`
	OnlyAllowMembersViewMembership bool   `json:"OnlyAllowMembersViewMembership"`
	OwnerTitle                     string `json:"OwnerTitle"`
	PrincipalType                  int    `json:"PrincipalType"`
	RequestToJoinLeaveEmailSetting bool   `json:"RequestToJoinLeaveEmailSetting"`
	Title                          string `json:"Title"`
}

// GroupResp - group response type with helper processor methods
type GroupResp []byte

// NewGroup - Group struct constructor function
func NewGroup(client *gosip.SPClient, endpoint string, config *RequestConfig) *Group {
	return &Group{
		client:    client,
		endpoint:  endpoint,
		config:    config,
		modifiers: NewODataMods(),
	}
}

// ToURL gets endpoint with modificators raw URL
func (group *Group) ToURL() string {
	return toURL(group.endpoint, group.modifiers)
}

// Get gets group data object
func (group *Group) Get() (GroupResp, error) {
	client := NewHTTPClient(group.client)
	return client.Get(group.ToURL(), getConfHeaders(group.config))
}

// Update updates Group's metadata with properties provided in `body` parameter
// where `body` is byte array representation of JSON string payload relevalt to SP.Group object
func (group *Group) Update(body []byte) (GroupResp, error) {
	body = patchMetadataType(body, "SP.Group")
	client := NewHTTPClient(group.client)
	return client.Update(group.endpoint, bytes.NewBuffer(body), getConfHeaders(group.config))
}

// Users gets Users API queryable collection
func (group *Group) Users() *Users {
	return NewUsers(
		group.client,
		fmt.Sprintf("%s/Users", group.endpoint),
		group.config,
	)
}

// AddUser adds a user by login name to this group
func (group *Group) AddUser(loginName string) error {
	endpoint := fmt.Sprintf("%s/Users", group.ToURL())
	client := NewHTTPClient(group.client)
	metadata := make(map[string]interface{})
	metadata["__metadata"] = map[string]string{
		"type": "SP.User",
	}
	metadata["LoginName"] = loginName
	body, _ := json.Marshal(metadata)
	_, err := client.Post(endpoint, bytes.NewBuffer(body), getConfHeaders(group.config))
	return err
}

// AddUserByID adds a user by ID to this group
func (group *Group) AddUserByID(userID int) error {
	sp := NewSP(group.client)
	user, err := sp.Web().SiteUsers().Select("LoginName").GetByID(userID).Get()
	if err != nil {
		return err
	}
	return group.AddUser(user.Data().LoginName)
}

// SetAsOwner sets a user as owner
func (group *Group) SetAsOwner(userID int) error {
	endpoint := fmt.Sprintf("%s/SetUserAsOwner(%d)", group.ToURL(), userID)
	client := NewHTTPClient(group.client)
	_, err := client.Post(endpoint, nil, getConfHeaders(group.config))
	return err
}

// RemoveUser removes a user from group
func (group *Group) RemoveUser(loginName string) error {
	endpoint := fmt.Sprintf(
		"%s/Users/RemoveByLoginName(@v)?@v='%s'",
		group.ToURL(),
		url.QueryEscape(loginName),
	)
	client := NewHTTPClient(group.client)
	_, err := client.Post(endpoint, nil, getConfHeaders(group.config))
	return err
}

// RemoveUserByID removes a user from group
func (group *Group) RemoveUserByID(userID int) error {
	endpoint := fmt.Sprintf("%s/Users/RemoveById(%d)", group.ToURL(), userID)
	client := NewHTTPClient(group.client)
	_, err := client.Post(endpoint, nil, getConfHeaders(group.config))
	return err
}
