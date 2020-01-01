package api

import (
	"encoding/json"
	"fmt"

	"github.com/koltyakov/gosip"
)

// ContentType represents SharePoint Content Types API queryable object struct
type ContentType struct {
	client    *gosip.SPClient
	config    *RequestConfig
	endpoint  string
	modifiers map[string]string
}

// ContentTypeInfo - content type API response payload structure
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

// ContentTypeResp - content type response type with helper processor methods
type ContentTypeResp []byte

// NewContentType - ContentType struct constructor function
func NewContentType(client *gosip.SPClient, endpoint string, config *RequestConfig) *ContentType {
	return &ContentType{
		client:   client,
		endpoint: endpoint,
		config:   config,
	}
}

// ToURL gets endpoint with modificators raw URL ...
func (ct *ContentType) ToURL() string {
	return toURL(ct.endpoint, ct.modifiers)
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
	return json.Unmarshal(data, obj)
}
