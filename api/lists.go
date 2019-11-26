package api

import (
	"fmt"
	"net/url"

	"github.com/koltyakov/gosip"
)

// Lists ...
type Lists struct {
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
func (lists *Lists) Conf(conf *Conf) *Lists {
	lists.conf = conf
	return lists
}

// Select ...
func (lists *Lists) Select(oDataSelect string) *Lists {
	lists.oSelect = oDataSelect
	return lists
}

// Expand ...
func (lists *Lists) Expand(oDataExpand string) *Lists {
	lists.oExpand = oDataExpand
	return lists
}

// Filter ...
func (lists *Lists) Filter(oDataFilter string) *Lists {
	lists.oFilter = oDataFilter
	return lists
}

// Top ...
func (lists *Lists) Top(oDataTop int) *Lists {
	lists.oTop = oDataTop
	return lists
}

// OrderBy ...
func (lists *Lists) OrderBy(oDataOrderBy string, ascending bool) *Lists {
	direction := "asc"
	if !ascending {
		direction = "desc"
	}
	if lists.oOrderBy != "" {
		lists.oOrderBy += ","
	}
	lists.oOrderBy += fmt.Sprintf("%s %s", oDataOrderBy, direction)
	return lists
}

// Get ...
func (lists *Lists) Get() ([]byte, error) {
	apiURL, _ := url.Parse(lists.endpoint)
	query := url.Values{}
	if lists.oSelect != "" {
		query.Add("$select", TrimMultiline(lists.oSelect))
	}
	if lists.oExpand != "" {
		query.Add("$expand", TrimMultiline(lists.oExpand))
	}
	if lists.oFilter != "" {
		query.Add("$filter", TrimMultiline(lists.oFilter))
	}
	if lists.oTop != 0 {
		query.Add("$top", fmt.Sprintf("%d", lists.oTop))
	}
	if lists.oOrderBy != "" {
		query.Add("$orderBy", TrimMultiline(lists.oOrderBy))
	}
	apiURL.RawQuery = query.Encode()
	sp := &HTTPClient{SPClient: lists.client}
	headers := map[string]string{}
	if lists.conf != nil {
		headers = lists.conf.Headers
	}
	return sp.Get(apiURL.String(), headers)
}

// GetByTitle ...
func (lists *Lists) GetByTitle(listTitle string) *List {
	list := &List{
		client: lists.client,
		conf:   lists.conf,
		endpoint: fmt.Sprintf(
			"%s/GetByTitle('%s')",
			lists.endpoint,
			listTitle,
		),
	}
	return list
}

// GetByID ...
func (lists *Lists) GetByID(listGUID string) *List {
	list := &List{
		client: lists.client,
		conf:   lists.conf,
		endpoint: fmt.Sprintf(
			"%s('%s')",
			lists.endpoint,
			listGUID,
		),
	}
	return list
}
