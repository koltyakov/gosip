package api

import (
	"fmt"
	"net/url"

	"github.com/koltyakov/gosip"
)

// Lists ...
type Lists struct {
	client    *gosip.SPClient
	config    *RequestConfig
	endpoint  string
	modifiers map[string]string
}

// NewLists ...
func NewLists(client *gosip.SPClient, endpoint string, config *RequestConfig) *Lists {
	return &Lists{
		client:   client,
		endpoint: endpoint,
		config:   config,
	}
}

// ToURL ...
func (lists *Lists) ToURL() string {
	return lists.endpoint
}

// Conf ...
func (lists *Lists) Conf(config *RequestConfig) *Lists {
	lists.config = config
	return lists
}

// Select ...
func (lists *Lists) Select(oDataSelect string) *Lists {
	if lists.modifiers == nil {
		lists.modifiers = make(map[string]string)
	}
	lists.modifiers["$select"] = oDataSelect
	return lists
}

// Expand ...
func (lists *Lists) Expand(oDataExpand string) *Lists {
	if lists.modifiers == nil {
		lists.modifiers = make(map[string]string)
	}
	lists.modifiers["$expand"] = oDataExpand
	return lists
}

// Filter ...
func (lists *Lists) Filter(oDataFilter string) *Lists {
	if lists.modifiers == nil {
		lists.modifiers = make(map[string]string)
	}
	lists.modifiers["$filter"] = oDataFilter
	return lists
}

// Top ...
func (lists *Lists) Top(oDataTop int) *Lists {
	if lists.modifiers == nil {
		lists.modifiers = make(map[string]string)
	}
	lists.modifiers["$top"] = string(oDataTop)
	return lists
}

// OrderBy ...
func (lists *Lists) OrderBy(oDataOrderBy string, ascending bool) *Lists {
	direction := "asc"
	if !ascending {
		direction = "desc"
	}
	if lists.modifiers == nil {
		lists.modifiers = make(map[string]string)
	}
	if lists.modifiers["$orderby"] != "" {
		lists.modifiers["$orderby"] += ","
	}
	lists.modifiers["$orderby"] += fmt.Sprintf("%s %s", oDataOrderBy, direction)
	return lists
}

// Get ...
func (lists *Lists) Get() ([]byte, error) {
	apiURL, _ := url.Parse(lists.endpoint)
	query := url.Values{}
	for k, v := range lists.modifiers {
		query.Add(k, trimMultiline(v))
	}
	apiURL.RawQuery = query.Encode()
	sp := NewHTTPClient(lists.client)
	headers := map[string]string{}
	if lists.config != nil {
		headers = lists.config.Headers
	}
	return sp.Get(apiURL.String(), headers)
}

// GetByTitle ...
func (lists *Lists) GetByTitle(listTitle string) *List {
	list := NewList(
		lists.client,
		fmt.Sprintf("%s/GetByTitle('%s')", lists.endpoint, listTitle),
		lists.config,
	)
	return list
}

// GetByID ...
func (lists *Lists) GetByID(listGUID string) *List {
	list := NewList(
		lists.client,
		fmt.Sprintf("%s('%s')", lists.endpoint, listGUID),
		lists.config,
	)
	return list
}
