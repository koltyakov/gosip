package api

import (
	"fmt"
	"net/url"

	"github.com/koltyakov/gosip"
)

// Group ...
type Group struct {
	client   *gosip.SPClient
	config   *RequestConfig
	endpoint string
	oSelect  string
	oExpand  string
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

// Conf ...
func (group *Group) Conf(config *RequestConfig) *Group {
	group.config = config
	return group
}

// Select ...
func (group *Group) Select(oDataSelect string) *Group {
	group.oSelect = oDataSelect
	return group
}

// Expand ...
func (group *Group) Expand(oDataExpand string) *Group {
	group.oExpand = oDataExpand
	return group
}

// Get ...
func (group *Group) Get() ([]byte, error) {
	apiURL, _ := url.Parse(group.endpoint)
	query := url.Values{}
	if group.oSelect != "" {
		query.Add("$select", trimMultiline(group.oSelect))
	}
	if group.oExpand != "" {
		query.Add("$expand", trimMultiline(group.oExpand))
	}
	apiURL.RawQuery = query.Encode()
	sp := NewHTTPClient(group.client)
	return sp.Get(apiURL.String(), getConfHeaders(group.config))
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
	return &Users{
		client: group.client,
		config: group.config,
		endpoint: fmt.Sprintf(
			"%s/Users",
			group.endpoint,
		),
	}
}
