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

// ContentTypeResp ...
type ContentTypeResp []byte

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
func (ct *ContentType) Get() (ContentTypeResp, error) {
	sp := NewHTTPClient(ct.client)
	data, err := sp.Get(ct.ToURL(), getConfHeaders(ct.config))
	if err != nil {
		return nil, err
	}
	return ContentTypeResp(data), nil
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

/* Response helpers */

// Data : to get typed data
func (ctResp *ContentTypeResp) Data() *ContentTypeInfo {
	data := parseODataItem(*ctResp)
	res := &ContentTypeInfo{}
	json.Unmarshal(data, &res)
	return res
}

// Unmarshal : to unmarshal to custom object
func (ctResp *ContentTypeResp) Unmarshal(obj interface{}) error {
	data := parseODataItem(*ctResp)
	return json.Unmarshal(data, &obj)
}
