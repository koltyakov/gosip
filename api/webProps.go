package api

import (
	"net/url"

	"github.com/koltyakov/gosip"
)

// WebProps ...
type WebProps struct {
	client    *gosip.SPClient
	config    *RequestConfig
	endpoint  string
	modifiers map[string]string
}

// NewWebProps ...
func NewWebProps(client *gosip.SPClient, endpoint string, config *RequestConfig) *WebProps {
	return &WebProps{
		client:   client,
		endpoint: endpoint,
		config:   config,
	}
}

// ToURL ...
func (webProps *WebProps) ToURL() string {
	apiURL, _ := url.Parse(webProps.endpoint)
	query := url.Values{}
	for k, v := range webProps.modifiers {
		query.Add(k, trimMultiline(v))
	}
	apiURL.RawQuery = query.Encode()
	return apiURL.String()
}

// Conf ...
func (webProps *WebProps) Conf(config *RequestConfig) *WebProps {
	webProps.config = config
	return webProps
}

// Select ...
func (webProps *WebProps) Select(oDataSelect string) *WebProps {
	if webProps.modifiers == nil {
		webProps.modifiers = make(map[string]string)
	}
	webProps.modifiers["$select"] = oDataSelect
	return webProps
}

// Expand ...
func (webProps *WebProps) Expand(oDataExpand string) *WebProps {
	if webProps.modifiers == nil {
		webProps.modifiers = make(map[string]string)
	}
	webProps.modifiers["$expand"] = oDataExpand
	return webProps
}

// Get ...
func (webProps *WebProps) Get() ([]byte, error) {
	sp := NewHTTPClient(webProps.client)
	headers := map[string]string{}
	if webProps.config != nil {
		headers = webProps.config.Headers
	}
	return sp.Get(webProps.ToURL(), headers)
}

// ToDo:
// Write Props with CSOM
