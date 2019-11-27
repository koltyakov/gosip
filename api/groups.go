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
	groups.modifiers["$top"] = string(oDataTop)
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
	apiURL, _ := url.Parse(groups.endpoint)
	query := url.Values{}
	for k, v := range groups.modifiers {
		query.Add(k, trimMultiline(v))
	}
	apiURL.RawQuery = query.Encode()
	sp := NewHTTPClient(groups.client)
	return sp.Get(apiURL.String(), getConfHeaders(groups.config))
}

// Add ...
func (groups *Groups) Add(body []byte) ([]byte, error) {
	sp := NewHTTPClient(groups.client)
	return sp.Post(groups.endpoint, body, getConfHeaders(groups.config))
}

// GetByID ...
func (groups *Groups) GetByID(groupID int) *Group {
	return &Group{
		client: groups.client,
		config: groups.config,
		endpoint: fmt.Sprintf("%s/GetById(%d)",
			groups.endpoint,
			groupID,
		),
	}
}

// GetByName ...
func (groups *Groups) GetByName(groupName string) *Group {
	return &Group{
		client: groups.client,
		config: groups.config,
		endpoint: fmt.Sprintf("%s/GetByName('%s')",
			groups.endpoint,
			groupName,
		),
	}
}