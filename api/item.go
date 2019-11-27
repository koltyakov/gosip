package api

import (
	"net/url"

	"github.com/koltyakov/gosip"
)

// Item ...
type Item struct {
	client   *gosip.SPClient
	config   *RequestConfig
	endpoint string
	oSelect  string
	oExpand  string
}

// Conf ...
func (item *Item) Conf(config *RequestConfig) *Item {
	item.config = config
	return item
}

// Select ...
func (item *Item) Select(oDataSelect string) *Item {
	item.oSelect = oDataSelect
	return item
}

// Expand ...
func (item *Item) Expand(oDataExpand string) *Item {
	item.oExpand = oDataExpand
	return item
}

// Get ...
func (item *Item) Get() ([]byte, error) {
	apiURL, _ := url.Parse(item.endpoint)
	query := url.Values{}
	if item.oSelect != "" {
		query.Add("$select", trimMultiline(item.oSelect))
	}
	if item.oExpand != "" {
		query.Add("$expand", trimMultiline(item.oExpand))
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
