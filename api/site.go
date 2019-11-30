package api

import (
	"fmt"
	"net/url"

	"github.com/koltyakov/gosip"
)

// Site ...
type Site struct {
	client    *gosip.SPClient
	config    *RequestConfig
	endpoint  string
	modifiers map[string]string
}

// NewSite ...
func NewSite(client *gosip.SPClient, endpoint string, config *RequestConfig) *Site {
	return &Site{
		client:   client,
		endpoint: endpoint,
		config:   config,
	}
}

// ToURL ...
func (site *Site) ToURL() string {
	apiURL, _ := url.Parse(site.endpoint)
	query := url.Values{}
	for k, v := range site.modifiers {
		query.Add(k, trimMultiline(v))
	}
	apiURL.RawQuery = query.Encode()
	return apiURL.String()
}

// Conf ...
func (site *Site) Conf(config *RequestConfig) *Site {
	site.config = config
	return site
}

// Select ...
func (site *Site) Select(oDataSelect string) *Site {
	if site.modifiers == nil {
		site.modifiers = make(map[string]string)
	}
	site.modifiers["$select"] = oDataSelect
	return site
}

// Expand ...
func (site *Site) Expand(oDataExpand string) *Site {
	if site.modifiers == nil {
		site.modifiers = make(map[string]string)
	}
	site.modifiers["$expand"] = oDataExpand
	return site
}

// Get ...
func (site *Site) Get() ([]byte, error) {
	sp := NewHTTPClient(site.client)
	return sp.Get(site.ToURL(), getConfHeaders(site.config))
}

// Delete ...
func (site *Site) Delete() ([]byte, error) {
	sp := NewHTTPClient(site.client)
	return sp.Delete(site.endpoint, getConfHeaders(site.config))
}

// RootWeb ...
func (site *Site) RootWeb() *Web {
	endpoint := fmt.Sprintf("%s/RootWeb", site.endpoint)
	return NewWeb(site.client, endpoint, site.config)
}

// OpenWebByID ...
func (site *Site) OpenWebByID(webID string) ([]byte, error) {
	endpoint := fmt.Sprintf("%s/OpenWebById('%s')", site.endpoint, webID)
	sp := NewHTTPClient(site.client)
	return sp.Post(endpoint, nil, getConfHeaders(site.config))
}

// ToDo:
// Features
// Custom actions
