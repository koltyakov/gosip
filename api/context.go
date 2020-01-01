package api

import (
	"encoding/json"
	"fmt"

	"github.com/koltyakov/gosip"
)

// Context represents SharePoint Content Info API object struct
// Always use NewContext constructor instead of &Context{}
type Context struct {
	client   *gosip.SPClient
	config   *RequestConfig
	endpoint string
}

// ContextInfo - context info response payload structure
type ContextInfo struct {
	FormDigestTimeoutSeconds int      `json:"FormDigestTimeoutSeconds"`
	FormDigestValue          string   `json:"FormDigestValue"`
	LibraryVersion           string   `json:"LibraryVersion"` // "16.0.19520.12058",
	SiteFullURL              string   `json:"SiteFullUrl"`
	SupportedSchemaVersions  []string `json:"SupportedSchemaVersions"`
	WebFullURL               string   `json:"WebFullUrl"`
}

// NewContext - Context struct constructor function
func NewContext(client *gosip.SPClient, endpoint string, config *RequestConfig) *Context {
	return &Context{
		client:   client,
		endpoint: endpoint,
		config:   config,
	}
}

// Get gets context info data object
func (context *Context) Get() (*ContextInfo, error) {
	endpoint := fmt.Sprintf("%s/_api/ContextInfo", getPriorEndpoint(context.endpoint, "/_api"))
	resp, err := NewHTTPClient(context.client).Post(endpoint, nil, getConfHeaders(context.config))
	if err != nil {
		return nil, err
	}
	data := parseODataItem(resp)
	ci0 := &struct {
		GetContextWebInformation map[string]interface{} `json:"GetContextWebInformation"`
	}{}
	if err := json.Unmarshal(data, &ci0); err != nil {
		return nil, err
	}
	ci0.GetContextWebInformation = normalizeMultiLookupsMap(ci0.GetContextWebInformation)
	data1, err := json.Marshal(ci0.GetContextWebInformation)
	if err != nil {
		return nil, err
	}
	contextInfo := &ContextInfo{}
	if err := json.Unmarshal(data1, &contextInfo); err != nil {
		return nil, err
	}
	return contextInfo, err
}
