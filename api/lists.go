package api

import (
	"encoding/json"
	"fmt"

	"github.com/koltyakov/gosip"
)

// Lists represent SharePoint Lists API queryable collection struct
// Always use NewLists constructor instead of &Lists{}
type Lists struct {
	client    *gosip.SPClient
	config    *RequestConfig
	endpoint  string
	modifiers *ODataMods
}

// ListsResp - lists response type with helper processor methods
type ListsResp []byte

// NewLists - Lists struct constructor function
func NewLists(client *gosip.SPClient, endpoint string, config *RequestConfig) *Lists {
	return &Lists{
		client:    client,
		endpoint:  endpoint,
		config:    config,
		modifiers: NewODataMods(),
	}
}

// ToURL gets endpoint with modificators raw URL
func (lists *Lists) ToURL() string {
	return toURL(lists.endpoint, lists.modifiers)
}

// Conf receives custom request config definition, e.g. custom headers, custom OData mod
func (lists *Lists) Conf(config *RequestConfig) *Lists {
	lists.config = config
	return lists
}

// Select ...
func (lists *Lists) Select(oDataSelect string) *Lists {
	lists.modifiers.AddSelect(oDataSelect)
	return lists
}

// Expand ...
func (lists *Lists) Expand(oDataExpand string) *Lists {
	lists.modifiers.AddExpand(oDataExpand)
	return lists
}

// Filter ...
func (lists *Lists) Filter(oDataFilter string) *Lists {
	lists.modifiers.AddFilter(oDataFilter)
	return lists
}

// Top ...
func (lists *Lists) Top(oDataTop int) *Lists {
	lists.modifiers.AddTop(oDataTop)
	return lists
}

// OrderBy ...
func (lists *Lists) OrderBy(oDataOrderBy string, ascending bool) *Lists {
	lists.modifiers.AddOrderBy(oDataOrderBy, ascending)
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
	collection, _ := parseODataCollection(*listsResp)
	lists := []ListResp{}
	for _, list := range collection {
		lists = append(lists, ListResp(list))
	}
	return lists
}

// Unmarshal : to unmarshal to custom object
func (listsResp *ListsResp) Unmarshal(obj interface{}) error {
	// collection := parseODataCollection(*listsResp)
	// data, _ := json.Marshal(collection)
	data, _ := parseODataCollectionPlain(*listsResp)
	return json.Unmarshal(data, obj)
}
