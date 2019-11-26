package api

import (
	"net/url"

	"github.com/koltyakov/gosip"
)

// Group ...
type Group struct {
	client   *gosip.SPClient
	conf     *Conf
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
func (group *Group) Conf(conf *Conf) *Group {
	group.conf = conf
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
		query.Add("$select", TrimMultiline(group.oSelect))
	}
	if group.oExpand != "" {
		query.Add("$expand", TrimMultiline(group.oExpand))
	}
	apiURL.RawQuery = query.Encode()
	sp := &HTTPClient{SPClient: group.client}
	return sp.Get(apiURL.String(), GetConfHeaders(group.conf))
}

// Delete ...
func (group *Group) Delete() ([]byte, error) {
	sp := &HTTPClient{SPClient: group.client}
	return sp.Delete(group.endpoint, GetConfHeaders(group.conf))
}

// Update ...
func (group *Group) Update(body []byte) ([]byte, error) {
	sp := &HTTPClient{SPClient: group.client}
	return sp.Update(group.endpoint, body, GetConfHeaders(group.conf))
}
