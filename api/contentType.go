package api

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/koltyakov/gosip"
)

// ContentType ...
type ContentType struct {
	client    *gosip.SPClient
	config    *RequestConfig
	endpoint  string
	modifiers map[string]string
}

// ContentTypeInfo ...
type ContentTypeInfo struct {
	Description string `json:"Description"`
	Group       string `json:"Group"`
	Hidden      bool   `json:"Hidden"`
	JSLink      string `json:"JSLink"`
	Name        string `json:"Name"`
	ReadOnly    bool   `json:"Read"`
	SchemaXML   string `json:"SchemaXml"`
	Scope       string `json:"Scope"`
	Sealed      bool   `json:"Sealed"`
	ID          string `json:"StringId"`
}

// NewContentType ...
func NewContentType(client *gosip.SPClient, endpoint string, config *RequestConfig) *ContentType {
	return &ContentType{
		client:   client,
		endpoint: endpoint,
		config:   config,
	}
}

// ToURL ...
func (ct *ContentType) ToURL() string {
	apiURL, _ := url.Parse(ct.endpoint)
	query := url.Values{}
	for k, v := range ct.modifiers {
		query.Add(k, trimMultiline(v))
	}
	apiURL.RawQuery = query.Encode()
	return apiURL.String()
}

// Conf ...
func (ct *ContentType) Conf(config *RequestConfig) *ContentType {
	ct.config = config
	return ct
}

// Select ...
func (ct *ContentType) Select(oDataSelect string) *ContentType {
	if ct.modifiers == nil {
		ct.modifiers = make(map[string]string)
	}
	ct.modifiers["$select"] = oDataSelect
	return ct
}

// Expand ...
func (ct *ContentType) Expand(oDataExpand string) *ContentType {
	if ct.modifiers == nil {
		ct.modifiers = make(map[string]string)
	}
	ct.modifiers["$expand"] = oDataExpand
	return ct
}

// Get ...
func (ct *ContentType) Get() (*ContentTypeInfo, error) {
	sp := NewHTTPClient(ct.client)
	data, err := sp.Get(ct.ToURL(), HeadersPresets.Verbose.Headers)
	if err != nil {
		return nil, err
	}
	res := &struct {
		CT *ContentTypeInfo `json:"d"`
	}{}
	if err := json.Unmarshal(data, &res); err != nil {
		return nil, err
	}
	return res.CT, nil
}

// GetRaw ...
func (ct *ContentType) GetRaw() ([]byte, error) {
	sp := NewHTTPClient(ct.client)
	return sp.Get(ct.ToURL(), getConfHeaders(ct.config))
}

// Delete ...
func (ct *ContentType) Delete() ([]byte, error) {
	sp := NewHTTPClient(ct.client)
	return sp.Delete(ct.endpoint, getConfHeaders(ct.config))
}

// Recycle ...
func (ct *ContentType) Recycle() ([]byte, error) {
	sp := NewHTTPClient(ct.client)
	endpoint := fmt.Sprintf("%s/Recycle", ct.endpoint)
	return sp.Post(endpoint, nil, getConfHeaders(ct.config))
}

// Fields ...
func (ct *ContentType) Fields() *Fields {
	return NewFields(
		ct.client,
		fmt.Sprintf("%s/Fields", ct.endpoint),
		ct.config,
	)
}
