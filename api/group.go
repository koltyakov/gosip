package api

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/koltyakov/gosip"
)

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

// Conf receives custom request config definition, e.g. custom headers, custom OData mod
func (group *Group) Conf(config *RequestConfig) *Group {
	group.config = config
	return group
}

// Select adds $select OData modifier
func (group *Group) Select(oDataSelect string) *Group {
	group.modifiers.AddSelect(oDataSelect)
	return group
}

// Expand adds $expand OData modifier
func (group *Group) Expand(oDataExpand string) *Group {
	group.modifiers.AddExpand(oDataExpand)
	return group
}

// Get gets group data object
func (group *Group) Get() (GroupResp, error) {
	sp := NewHTTPClient(group.client)
	return sp.Get(group.ToURL(), getConfHeaders(group.config))
}

// Update updates Group's metadata with properties provided in `body` parameter
// where `body` is byte array representation of JSON string payload relevalt to SP.Group object
func (group *Group) Update(body []byte) (GroupResp, error) {
	body = patchMetadataType(body, "SP.Group")
	sp := NewHTTPClient(group.client)
	return sp.Update(group.endpoint, body, getConfHeaders(group.config))
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
	sp := NewHTTPClient(group.client)
	metadata := make(map[string]interface{})
	metadata["__metadata"] = map[string]string{
		"type": "SP.User",
	}
	metadata["LoginName"] = loginName
	body, _ := json.Marshal(metadata)
	_, err := sp.Post(endpoint, body, getConfHeaders(group.config))
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
	sp := NewHTTPClient(group.client)
	_, err := sp.Post(endpoint, nil, getConfHeaders(group.config))
	return err
}

// RemoveUser removes a user from group
func (group *Group) RemoveUser(loginName string) error {
	endpoint := fmt.Sprintf(
		"%s/Users/RemoveByLoginName(@v)?@v='%s'",
		group.ToURL(),
		url.QueryEscape(loginName),
	)
	sp := NewHTTPClient(group.client)
	_, err := sp.Post(endpoint, nil, getConfHeaders(group.config))
	return err
}

// RemoveUserByID removes a user from group
func (group *Group) RemoveUserByID(userID int) error {
	endpoint := fmt.Sprintf("%s/Users/RemoveById(%d)", group.ToURL(), userID)
	sp := NewHTTPClient(group.client)
	_, err := sp.Post(endpoint, nil, getConfHeaders(group.config))
	return err
}

/* Response helpers */

// Data : to get typed data
func (groupResp *GroupResp) Data() *GroupInfo {
	data := NormalizeODataItem(*groupResp)
	res := &GroupInfo{}
	json.Unmarshal(data, &res)
	return res
}

// Unmarshal : to unmarshal to custom object
func (groupResp *GroupResp) Unmarshal(obj interface{}) error {
	data := NormalizeODataItem(*groupResp)
	return json.Unmarshal(data, obj)
}
