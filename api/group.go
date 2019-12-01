package api

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/koltyakov/gosip"
)

// Group ...
type Group struct {
	client    *gosip.SPClient
	config    *RequestConfig
	endpoint  string
	modifiers map[string]string
}

// GroupInfo ...
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

// GroupResp ...
type GroupResp []byte

// NewGroup ...
func NewGroup(client *gosip.SPClient, endpoint string, config *RequestConfig) *Group {
	return &Group{
		client:   client,
		endpoint: endpoint,
		config:   config,
	}
}

// ToURL ...
func (group *Group) ToURL() string {
	apiURL, _ := url.Parse(group.endpoint)
	query := url.Values{}
	for k, v := range group.modifiers {
		query.Add(k, trimMultiline(v))
	}
	apiURL.RawQuery = query.Encode()
	return apiURL.String()
}

// Conf ...
func (group *Group) Conf(config *RequestConfig) *Group {
	group.config = config
	return group
}

// Select ...
func (group *Group) Select(oDataSelect string) *Group {
	if group.modifiers == nil {
		group.modifiers = make(map[string]string)
	}
	group.modifiers["$select"] = oDataSelect
	return group
}

// Expand ...
func (group *Group) Expand(oDataExpand string) *Group {
	if group.modifiers == nil {
		group.modifiers = make(map[string]string)
	}
	group.modifiers["$expand"] = oDataExpand
	return group
}

// Get ...
func (group *Group) Get() (GroupResp, error) {
	sp := NewHTTPClient(group.client)
	return sp.Get(group.ToURL(), getConfHeaders(group.config))
}

// Update ...
func (group *Group) Update(body []byte) ([]byte, error) {
	body = patchMetadataType(body, "SP.Group")
	sp := NewHTTPClient(group.client)
	return sp.Update(group.endpoint, body, getConfHeaders(group.config))
}

// Users ...
func (group *Group) Users() *Users {
	return NewUsers(
		group.client,
		fmt.Sprintf("%s/Users", group.endpoint),
		group.config,
	)
}

// AddUser ...
func (group *Group) AddUser(loginName string) ([]byte, error) {
	endpoint := fmt.Sprintf("%s/Users", group.ToURL())
	sp := NewHTTPClient(group.client)

	metadata := make(map[string]interface{})
	metadata["__metadata"] = map[string]string{
		"type": "SP.User",
	}
	metadata["LoginName"] = loginName
	body, _ := json.Marshal(metadata)
	return sp.Post(endpoint, body, getConfHeaders(group.config))
}

/* Response helpers */

// Data : to get typed data
func (groupResp *GroupResp) Data() *GroupInfo {
	data := parseODataItem(*groupResp)
	res := &GroupInfo{}
	json.Unmarshal(data, &res)
	return res
}

// Unmarshal : to unmarshal to custom object
func (groupResp *GroupResp) Unmarshal(obj *interface{}) error {
	return json.Unmarshal(*groupResp, &obj)
}
