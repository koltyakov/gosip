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

// Select adds $select OData modifier
func (lists *Lists) Select(oDataSelect string) *Lists {
	lists.modifiers.AddSelect(oDataSelect)
	return lists
}

// Expand adds $expand OData modifier
func (lists *Lists) Expand(oDataExpand string) *Lists {
	lists.modifiers.AddExpand(oDataExpand)
	return lists
}

// Filter adds $filter OData modifier
func (lists *Lists) Filter(oDataFilter string) *Lists {
	lists.modifiers.AddFilter(oDataFilter)
	return lists
}

// Top adds $top OData modifier
func (lists *Lists) Top(oDataTop int) *Lists {
	lists.modifiers.AddTop(oDataTop)
	return lists
}

// OrderBy adds $orderby OData modifier
func (lists *Lists) OrderBy(oDataOrderBy string, ascending bool) *Lists {
	lists.modifiers.AddOrderBy(oDataOrderBy, ascending)
	return lists
}

// Get gets Lists API queryable collection
func (lists *Lists) Get() (ListsResp, error) {
	sp := NewHTTPClient(lists.client)
	return sp.Get(lists.ToURL(), getConfHeaders(lists.config))
}

// GetByTitle gets a list by its Display Name (Title)
func (lists *Lists) GetByTitle(listTitle string) *List {
	list := NewList(
		lists.client,
		fmt.Sprintf("%s/GetByTitle('%s')", lists.endpoint, listTitle),
		lists.config,
	)
	return list
}

// GetByID gets a list by its ID (GUID)
func (lists *Lists) GetByID(listGUID string) *List {
	list := NewList(
		lists.client,
		fmt.Sprintf("%s('%s')", lists.endpoint, listGUID),
		lists.config,
	)
	return list
}

// Add creates new list on this web with a provided `title`.
// Along with title additional metadata can be provided in optional `metadata` string map object.
// `metadata` props should correspond to `SP.List` API type. Some props have defaults as BaseTemplate (100), AllowContentTypes (false), etc.
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

// AddWithURI creates new list on this web with a provided `title` and `uri`.
// `url` stands for a system friendly URI (e.g. `custom-list`) while `title` is a human friendly name (e.g. `Custom List`).
// Along with uri and title additional metadata can be provided in optional `metadata` string map object.
// `metadata` props should correspond to `SP.List` API type. Some props have defaults as BaseTemplate (100), AllowContentTypes (false), etc.
func (lists *Lists) AddWithURI(title string, uri string, metadata map[string]interface{}) ([]byte, error) {
	data, err := lists.Conf(HeadersPresets.Verbose).Add(uri, metadata)
	if err != nil {
		return nil, err
	}

	metadata = make(map[string]interface{})
	metadata["__metadata"] = map[string]string{"type": "SP.List"}
	metadata["Title"] = title
	body, _ := json.Marshal(metadata)

	return lists.GetByID(data.Data().ID).Update(body)
}

/* Response helpers */

// Data : to get typed data
func (listsResp *ListsResp) Data() []ListResp {
	collection, _ := normalizeODataCollection(*listsResp)
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
	data, _ := NormalizeODataCollection(*listsResp)
	return json.Unmarshal(data, obj)
}
