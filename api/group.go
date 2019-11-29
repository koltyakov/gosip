package api

import (
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
	Title                          int    `json:"Title"`
}

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
func (group *Group) Get() ([]byte, error) {
	sp := NewHTTPClient(group.client)
	return sp.Get(group.ToURL(), getConfHeaders(group.config))
}

// Delete ...
func (group *Group) Delete() ([]byte, error) {
	sp := NewHTTPClient(group.client)
	return sp.Delete(group.endpoint, getConfHeaders(group.config))
}

// Update ...
func (group *Group) Update(body []byte) ([]byte, error) {
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
