package api

import (
	"fmt"
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

// NewItem ...
func NewItem(client *gosip.SPClient, endpoint string, config *RequestConfig) *Item {
	return &Item{
		client:   client,
		endpoint: endpoint,
		config:   config,
	}
}

// ToURL ...
func (item *Item) ToURL() string {
	apiURL, _ := url.Parse(item.endpoint)
	query := url.Values{}
	for k, v := range item.modifiers {
		query.Add(k, trimMultiline(v))
	}
	apiURL.RawQuery = query.Encode()
	return apiURL.String()
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
	sp := NewHTTPClient(item.client)
	return sp.Get(item.ToURL(), getConfHeaders(item.config))
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

// Recycle ...
func (item *Item) Recycle() ([]byte, error) {
	endpoint := fmt.Sprintf("%s/Recycle", item.endpoint)
	sp := NewHTTPClient(item.client)
	return sp.Post(endpoint, nil, getConfHeaders(item.config))
}

// Roles ...
func (item *Item) Roles() *Roles {
	return NewRoles(item.client, item.endpoint, item.config)
}
