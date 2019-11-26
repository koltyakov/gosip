package api

import (
	"net/url"

	"github.com/koltyakov/gosip"
)

// Item ...
type Item struct {
	client   *gosip.SPClient
	conf     *Conf
	endpoint string
	oSelect  string
	oExpand  string
}

// Conf ...
func (item *Item) Conf(conf *Conf) *Item {
	item.conf = conf
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
		query.Add("$select", TrimMultiline(item.oSelect))
	}
	if item.oExpand != "" {
		query.Add("$expand", TrimMultiline(item.oExpand))
	}
	apiURL.RawQuery = query.Encode()
	sp := &HTTPClient{SPClient: item.client}
	return sp.Get(apiURL.String(), GetConfHeaders(item.conf))
}

// Delete ...
func (item *Item) Delete() ([]byte, error) {
	sp := &HTTPClient{SPClient: item.client}
	return sp.Delete(item.endpoint, GetConfHeaders(item.conf))
}

// Update ...
func (item *Item) Update(body []byte) ([]byte, error) {
	sp := &HTTPClient{SPClient: item.client}
	return sp.Update(item.endpoint, body, GetConfHeaders(item.conf))
}
