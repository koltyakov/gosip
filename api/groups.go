package api

import (
	"fmt"
	"net/url"

	"github.com/koltyakov/gosip"
)

// Groups ...
type Groups struct {
	client    *gosip.SPClient
	config    *RequestConfig
	endpoint  string
	modifiers map[string]string
}

// NewGroups ...
func NewGroups(client *gosip.SPClient, endpoint string, config *RequestConfig) *Groups {
	return &Groups{
		client:   client,
		endpoint: endpoint,
		config:   config,
	}
}

// ToURL ...
func (groups *Groups) ToURL() string {
	apiURL, _ := url.Parse(groups.endpoint)
	query := url.Values{}
	for k, v := range groups.modifiers {
		query.Add(k, trimMultiline(v))
	}
	apiURL.RawQuery = query.Encode()
	return apiURL.String()
}

// Conf ...
func (groups *Groups) Conf(config *RequestConfig) *Groups {
	groups.config = config
	return groups
}

// Select ...
func (groups *Groups) Select(oDataSelect string) *Groups {
	if groups.modifiers == nil {
		groups.modifiers = make(map[string]string)
	}
	groups.modifiers["$select"] = oDataSelect
	return groups
}

// Expand ...
func (groups *Groups) Expand(oDataExpand string) *Groups {
	if groups.modifiers == nil {
		groups.modifiers = make(map[string]string)
	}
	groups.modifiers["$expand"] = oDataExpand
	return groups
}

// Filter ...
func (groups *Groups) Filter(oDataFilter string) *Groups {
	if groups.modifiers == nil {
		groups.modifiers = make(map[string]string)
	}
	groups.modifiers["$filter"] = oDataFilter
	return groups
}

// Top ...
func (groups *Groups) Top(oDataTop int) *Groups {
	if groups.modifiers == nil {
		groups.modifiers = make(map[string]string)
	}
	groups.modifiers["$top"] = fmt.Sprintf("%d", oDataTop)
	return groups
}

// OrderBy ...
func (groups *Groups) OrderBy(oDataOrderBy string, ascending bool) *Groups {
	direction := "asc"
	if !ascending {
		direction = "desc"
	}
	if groups.modifiers == nil {
		groups.modifiers = make(map[string]string)
	}
	if groups.modifiers["$orderby"] != "" {
		groups.modifiers["$orderby"] += ","
	}
	groups.modifiers["$orderby"] += fmt.Sprintf("%s %s", oDataOrderBy, direction)
	return groups
}

// Get ...
func (groups *Groups) Get() ([]byte, error) {
	sp := NewHTTPClient(groups.client)
	return sp.Get(groups.ToURL(), getConfHeaders(groups.config))
}

// Add ...
func (groups *Groups) Add(body []byte) ([]byte, error) {
	sp := NewHTTPClient(groups.client)
	return sp.Post(groups.endpoint, body, getConfHeaders(groups.config))
}

// GetByID ...
func (groups *Groups) GetByID(groupID int) *Group {
	return NewGroup(
		groups.client,
		fmt.Sprintf("%s/GetById(%d)", groups.endpoint, groupID),
		groups.config,
	)
}

// GetByName ...
func (groups *Groups) GetByName(groupName string) *Group {
	return NewGroup(
		groups.client,
		fmt.Sprintf("%s/GetByName('%s')", groups.endpoint, groupName),
		groups.config,
	)
}
