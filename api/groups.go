package api

import (
	"fmt"
	"net/url"

	"github.com/koltyakov/gosip"
)

// Groups ...
type Groups struct {
	client   *gosip.SPClient
	config   *RequestConfig
	endpoint string
	oSelect  string
	oExpand  string
	oFilter  string
	oTop     int
	oOrderBy string
}

// Conf ...
func (groups *Groups) Conf(config *RequestConfig) *Groups {
	groups.config = config
	return groups
}

// Select ...
func (groups *Groups) Select(oDataSelect string) *Groups {
	groups.oSelect = oDataSelect
	return groups
}

// Expand ...
func (groups *Groups) Expand(oDataExpand string) *Groups {
	groups.oExpand = oDataExpand
	return groups
}

// Filter ...
func (groups *Groups) Filter(oDataFilter string) *Groups {
	groups.oFilter = oDataFilter
	return groups
}

// Top ...
func (groups *Groups) Top(oDataTop int) *Groups {
	groups.oTop = oDataTop
	return groups
}

// OrderBy ...
func (groups *Groups) OrderBy(oDataOrderBy string, ascending bool) *Groups {
	direction := "asc"
	if !ascending {
		direction = "desc"
	}
	if groups.oOrderBy != "" {
		groups.oOrderBy += ","
	}
	groups.oOrderBy += fmt.Sprintf("%s %s", oDataOrderBy, direction)
	return groups
}

// Get ...
func (groups *Groups) Get() ([]byte, error) {
	apiURL, _ := url.Parse(groups.endpoint)
	query := url.Values{}
	if groups.oSelect != "" {
		query.Add("$select", trimMultiline(groups.oSelect))
	}
	if groups.oExpand != "" {
		query.Add("$expand", trimMultiline(groups.oExpand))
	}
	if groups.oFilter != "" {
		query.Add("$filter", trimMultiline(groups.oFilter))
	}
	if groups.oTop != 0 {
		query.Add("$top", fmt.Sprintf("%d", groups.oTop))
	}
	if groups.oOrderBy != "" {
		query.Add("$orderBy", trimMultiline(groups.oOrderBy))
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
