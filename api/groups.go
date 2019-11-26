package api

import (
	"fmt"
	"net/url"

	"github.com/koltyakov/gosip"
)

// Groups ...
type Groups struct {
	client   *gosip.SPClient
	conf     *Conf
	endpoint string
	oSelect  string
	oExpand  string
	oFilter  string
	oTop     int
	oOrderBy string
}

// Conf ...
func (groups *Groups) Conf(conf *Conf) *Groups {
	groups.conf = conf
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
		query.Add("$select", TrimMultiline(groups.oSelect))
	}
	if groups.oExpand != "" {
		query.Add("$expand", TrimMultiline(groups.oExpand))
	}
	if groups.oFilter != "" {
		query.Add("$filter", TrimMultiline(groups.oFilter))
	}
	if groups.oTop != 0 {
		query.Add("$top", fmt.Sprintf("%d", groups.oTop))
	}
	if groups.oOrderBy != "" {
		query.Add("$orderBy", TrimMultiline(groups.oOrderBy))
	}
	apiURL.RawQuery = query.Encode()
	sp := &HTTPClient{SPClient: groups.client}
	return sp.Get(apiURL.String(), GetConfHeaders(groups.conf))
}

// Add ...
func (groups *Groups) Add(body []byte) ([]byte, error) {
	sp := &HTTPClient{SPClient: groups.client}
	return sp.Post(groups.endpoint, body, GetConfHeaders(groups.conf))
}

// GetByID ...
func (groups *Groups) GetByID(groupID int) *Group {
	return &Group{
		client: groups.client,
		conf:   groups.conf,
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
		conf:   groups.conf,
		endpoint: fmt.Sprintf("%s/GetByName('%s')",
			groups.endpoint,
			groupName,
		),
	}
}
