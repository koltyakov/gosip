package api

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"github.com/koltyakov/gosip"
)

// Items ...
type Items struct {
	client    *gosip.SPClient
	config    *RequestConfig
	endpoint  string
	modifiers map[string]string
}

// ItemsResp ...
type ItemsResp []byte

// NewItems ...
func NewItems(client *gosip.SPClient, endpoint string, config *RequestConfig) *Items {
	return &Items{
		client:   client,
		endpoint: endpoint,
		config:   config,
	}
}

// ToURL ...
func (items *Items) ToURL() string {
	apiURL, _ := url.Parse(items.endpoint)
	query := apiURL.Query() // url.Values{}
	for k, v := range items.modifiers {
		query.Set(k, trimMultiline(v))
	}
	apiURL.RawQuery = query.Encode()
	return apiURL.String()
}

// Conf ...
func (items *Items) Conf(config *RequestConfig) *Items {
	items.config = config
	return items
}

// Select ...
func (items *Items) Select(oDataSelect string) *Items {
	if items.modifiers == nil {
		items.modifiers = make(map[string]string)
	}
	items.modifiers["$select"] = oDataSelect
	return items
}

// Expand ...
func (items *Items) Expand(oDataExpand string) *Items {
	if items.modifiers == nil {
		items.modifiers = make(map[string]string)
	}
	items.modifiers["$expand"] = oDataExpand
	return items
}

// Filter ...
func (items *Items) Filter(oDataFilter string) *Items {
	if items.modifiers == nil {
		items.modifiers = make(map[string]string)
	}
	items.modifiers["$filter"] = oDataFilter
	return items
}

// Top ...
func (items *Items) Top(oDataTop int) *Items {
	if items.modifiers == nil {
		items.modifiers = make(map[string]string)
	}
	items.modifiers["$top"] = fmt.Sprintf("%d", oDataTop)
	return items
}

// Skip ...
func (items *Items) Skip(skipToken string) *Items {
	if items.modifiers == nil {
		items.modifiers = make(map[string]string)
	}
	items.modifiers["$skiptoken"] = fmt.Sprintf("%s", skipToken)
	return items
}

// OrderBy ...
func (items *Items) OrderBy(oDataOrderBy string, ascending bool) *Items {
	direction := "asc"
	if !ascending {
		direction = "desc"
	}
	if items.modifiers == nil {
		items.modifiers = make(map[string]string)
	}
	if items.modifiers["$orderby"] != "" {
		items.modifiers["$orderby"] += ","
	}
	items.modifiers["$orderby"] += fmt.Sprintf("%s %s", oDataOrderBy, direction)
	return items
}

// Get ...
func (items *Items) Get() (ItemsResp, error) {
	sp := NewHTTPClient(items.client)
	return sp.Get(items.ToURL(), getConfHeaders(items.config))
}

// GetPaged ...
func (items *Items) GetPaged() (ItemsResp, func() (ItemsResp, error), error) {
	sp := NewHTTPClient(items.client)
	itemsResp, err := sp.Get(items.ToURL(), getConfHeaders(items.config))
	if err != nil {
		return nil, nil, err
	}
	getNextPage := func() (ItemsResp, error) {
		nextURL := getODataCollectionNextPageURL(itemsResp)
		if nextURL == "" {
			return nil, fmt.Errorf("unable to get next page")
		}
		return NewItems(items.client, nextURL, items.config).Get()
	}
	return itemsResp, getNextPage, nil
}

// GetAll ...
func (items *Items) GetAll() ([]ItemResp, error) {
	return getAll(nil, nil, items)
}

func getAll(res []ItemResp, cur ItemsResp, items *Items) ([]ItemResp, error) {
	if res == nil && cur == nil {
		itemsCopy := NewItems(items.client, items.endpoint, items.config)
		itemsCopy.modifiers = map[string]string{}
		for key, val := range items.modifiers {
			switch key {
			case "$select":
				itemsCopy.modifiers[key] = val
			case "$expand":
				itemsCopy.modifiers[key] = val
			case "$top":
				itemsCopy.modifiers[key] = val
			}
		}
		itemsResp, err := itemsCopy.Get()
		if err != nil {
			return nil, err
		}
		res = itemsResp.Data()
		cur = itemsResp
	}
	nextURL := getODataCollectionNextPageURL(cur)
	if nextURL == "" {
		return res, nil
	}
	nextItemsResp, err := NewItems(items.client, nextURL, items.config).Get()
	if err != nil {
		return res, err
	}
	res = append(res, nextItemsResp.Data()...)
	return getAll(res, nextItemsResp, items)
}

// Add ...
func (items *Items) Add(body []byte) (ItemResp, error) {
	body = patchMetadataTypeCB(body, func() string {
		endpoint := getPriorEndpoint(items.endpoint, "/Items")
		list := NewList(items.client, endpoint, nil)
		oDataType, _ := list.GetEntityType()
		return oDataType
	})
	sp := NewHTTPClient(items.client)
	return sp.Post(items.endpoint, body, getConfHeaders(items.config))
}

// GetByID ...
func (items *Items) GetByID(itemID int) *Item {
	return NewItem(
		items.client,
		fmt.Sprintf("%s(%d)", items.endpoint, itemID),
		items.config,
	)
}

// GetByCAML ...
func (items *Items) GetByCAML(caml string) (ItemsResp, error) {
	endpoint := fmt.Sprintf("%s/GetItems", strings.TrimRight(items.endpoint, "/Items"))
	apiURL, _ := url.Parse(endpoint)
	query := url.Values{}
	for k, v := range items.modifiers {
		query.Add(k, trimMultiline(v))
	}
	apiURL.RawQuery = query.Encode()

	body := trimMultiline(`{
		"query": {
			"__metadata": { "type": "SP.CamlQuery" },
			"ViewXml": "` + trimMultiline(caml) + `"
		}
	}`)
	sp := NewHTTPClient(items.client)
	return sp.Post(apiURL.String(), []byte(body), getConfHeaders(items.config))
}

// ToDo:
// Batch

/* Response helpers */

// Data : to get typed data
func (itemsResp *ItemsResp) Data() []ItemResp {
	collection, _ := parseODataCollection(*itemsResp)
	items := []ItemResp{}
	for _, item := range collection {
		items = append(items, ItemResp(item))
	}
	return items
}

// NextPageURL : gets next page OData collection
func (itemsResp *ItemsResp) NextPageURL() string {
	return getODataCollectionNextPageURL(*itemsResp)
}

// Unmarshal : to unmarshal to custom object
func (itemsResp *ItemsResp) Unmarshal(obj interface{}) error {
	// collection := parseODataCollection(*itemsResp)
	// data, _ := json.Marshal(collection)
	data, _ := parseODataCollectionPlain(*itemsResp)
	return json.Unmarshal(data, obj)
}
