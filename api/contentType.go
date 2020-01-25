package api

import (
	"encoding/json"
	"fmt"

	"github.com/koltyakov/gosip"
)

//go:generate ggen -ent ContentType -conf -mods Select,Expand

// ContentType represents SharePoint Content Types API queryable object struct
// Always use NewContentType constructor instead of &ContentType{}
type ContentType struct {
	client    *gosip.SPClient
	config    *RequestConfig
	endpoint  string
	modifiers *ODataMods
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
		client:    client,
		endpoint:  endpoint,
		config:    config,
		modifiers: NewODataMods(),
	}
}

// ToURL gets endpoint with modificators raw URL
func (ct *ContentType) ToURL() string {
	return toURL(ct.endpoint, ct.modifiers)
}

// Get gets content type data object
func (ct *ContentType) Get() (ContentTypeResp, error) {
	sp := NewHTTPClient(ct.client)
	return sp.Get(ct.ToURL(), getConfHeaders(ct.config))
}

// Update updates Content Types's metadata with properties provided in `body` parameter
// where `body` is byte array representation of JSON string payload relevalt to SP.ContentType object
func (ct *ContentType) Update(body []byte) (ContentTypeResp, error) {
	body = patchMetadataType(body, "SP.ContentType")
	sp := NewHTTPClient(ct.client)
	return sp.Update(ct.endpoint, body, getConfHeaders(ct.config))
}

// Delete deletes a content type skipping recycle bin
func (ct *ContentType) Delete() error {
	sp := NewHTTPClient(ct.client)
	_, err := sp.Delete(ct.endpoint, getConfHeaders(ct.config))
	return err
}

// FieldLinks gets FieldLinks API instance queryable collection
func (ct *ContentType) FieldLinks() *FieldLinks {
	return NewFieldLinks(
		ct.client,
		fmt.Sprintf("%s/FieldLinks", ct.endpoint),
		ct.config,
	)
}

/* Response helpers */

// Data : to get typed data
func (ctResp *ContentTypeResp) Data() *ContentTypeInfo {
	data := NormalizeODataItem(*ctResp)
	res := &ContentTypeInfo{}
	json.Unmarshal(data, &res)
	return res
}

// Normalized returns normalized body
func (ctResp *ContentTypeResp) Normalized() []byte {
	return NormalizeODataItem(*ctResp)
}
