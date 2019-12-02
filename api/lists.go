package api

import (
	"encoding/json"
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

// ListsResp ...
type ListsResp []byte

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
	apiURL, _ := url.Parse(lists.endpoint)
	query := url.Values{}
	for k, v := range lists.modifiers {
		query.Add(k, trimMultiline(v))
	}
	apiURL.RawQuery = query.Encode()
	return apiURL.String()
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
	lists.modifiers["$top"] = fmt.Sprintf("%d", oDataTop)
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
func (lists *Lists) Get() (ListsResp, error) {
	sp := NewHTTPClient(lists.client)
	return sp.Get(lists.ToURL(), getConfHeaders(lists.config))
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

// Add ...
func (lists *Lists) Add(title string, metadata map[string]interface{}) (ListResp, error) {
	if metadata == nil {
		metadata = make(map[string]interface{})
	}

	metadata["__metadata"] = map[string]string{"type": "SP.List"}

	metadata["Title"] = title

	// Default values
	if metadata["BaseTemplate"] == nil {
		metadata["BaseTemplate"] = 100
	}
	if metadata["AllowContentTypes"] == nil {
		metadata["AllowContentTypes"] = false
	}
	if metadata["ContentTypesEnabled"] == nil {
		metadata["ContentTypesEnabled"] = false
	}

	parameters, _ := json.Marshal(metadata)
	body := fmt.Sprintf("%s", parameters)

	sp := NewHTTPClient(lists.client)
	headers := getConfHeaders(lists.config)

	headers["Accept"] = "application/json;odata=verbose"
	headers["Content-Type"] = "application/json;odata=verbose;charset=utf-8"

	return sp.Post(lists.endpoint, []byte(body), headers)
}

// AddWithURI ...
func (lists *Lists) AddWithURI(title string, uri string, metadata map[string]interface{}) ([]byte, error) {
	data, err := lists.Conf(HeadersPresets.Verbose).Add(uri, metadata)
	if err != nil {
		return nil, err
	}

	res := &struct {
		D struct {
			ID string `json:"Id"`
		} `json:"d"`
	}{}

	if err := json.Unmarshal(data, &res); err != nil {
		return nil, err
	}

	metadata = make(map[string]interface{})
	metadata["__metadata"] = map[string]string{"type": "SP.List"}
	metadata["Title"] = title

	body, _ := json.Marshal(metadata)

	return lists.GetByID(res.D.ID).Update(body)
}

/* Response helpers */

// Data : to get typed data
func (listsResp *ListsResp) Data() []ListResp {
	collection := parseODataCollection(*listsResp)
	lists := []ListResp{}
	for _, list := range collection {
		lists = append(lists, ListResp(list))
	}
	return lists
}

// Unmarshal : to unmarshal to custom object
func (listsResp *ListsResp) Unmarshal(obj interface{}) error {
	collection := parseODataCollection(*listsResp)
	data, _ := json.Marshal(collection)
	return json.Unmarshal(data, &obj)
}
