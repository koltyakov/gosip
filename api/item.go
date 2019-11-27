package api

import (
	"net/url"

	"github.com/koltyakov/gosip"
)

// Item ...
type Item struct {
	client    *gosip.SPClient
	config    *RequestConfig
	endpoint  string
	modifiers map[string]string
}

// ToURL ...
func (item *Item) ToURL() string {
	return item.endpoint
}

// Conf ...
func (item *Item) Conf(config *RequestConfig) *Item {
	item.config = config
	return item
}

// Select ...
func (item *Item) Select(oDataSelect string) *Item {
	if item.modifiers == nil {
		item.modifiers = make(map[string]string)
	}
	item.modifiers["$select"] = oDataSelect
	return item
}

// Expand ...
func (item *Item) Expand(oDataExpand string) *Item {
	if item.modifiers == nil {
		item.modifiers = make(map[string]string)
	}
	item.modifiers["$expand"] = oDataExpand
	return item
}

// Get ...
func (item *Item) Get() ([]byte, error) {
	apiURL, _ := url.Parse(item.endpoint)
	query := url.Values{}
	for k, v := range item.modifiers {
		query.Add(k, trimMultiline(v))
	}
	apiURL.RawQuery = query.Encode()
	sp := NewHTTPClient(item.client)
	return sp.Get(apiURL.String(), getConfHeaders(item.config))
}

// Delete ...
func (item *Item) Delete() ([]byte, error) {
	sp := NewHTTPClient(item.client)
	return sp.Delete(item.endpoint, getConfHeaders(item.config))
}

// Update ...
func (item *Item) Update(body []byte) ([]byte, error) {
	sp := NewHTTPClient(item.client)
	return sp.Update(item.endpoint, body, getConfHeaders(item.config))
}

// Roles ...
func (item *Item) Roles() *Roles {
	return &Roles{
		client:   item.client,
		config:   item.config,
		endpoint: item.endpoint,
	}
}
