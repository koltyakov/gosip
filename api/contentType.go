package api

import (
	"bytes"
	"fmt"

	"github.com/koltyakov/gosip"
)

//go:generate ggen -ent ContentType -conf -mods Select,Expand -helpers Data,Normalized

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
func (contentType *ContentType) ToURL() string {
	return toURL(contentType.endpoint, contentType.modifiers)
}

// Get gets content type data object
func (contentType *ContentType) Get() (ContentTypeResp, error) {
	client := NewHTTPClient(contentType.client)
	return client.Get(contentType.ToURL(), contentType.config)
}

// Update updates Content Types's metadata with properties provided in `body` parameter
// where `body` is byte array representation of JSON string payload relevant to SP.ContentType object
func (contentType *ContentType) Update(body []byte) (ContentTypeResp, error) {
	body = patchMetadataType(body, "SP.ContentType")
	client := NewHTTPClient(contentType.client)
	return client.Update(contentType.endpoint, bytes.NewBuffer(body), contentType.config)
}

// Delete deletes a content type skipping recycle bin
func (contentType *ContentType) Delete() error {
	client := NewHTTPClient(contentType.client)
	_, err := client.Delete(contentType.endpoint, contentType.config)
	return err
}

// FieldLinks gets FieldLinks API instance queryable collection
func (contentType *ContentType) FieldLinks() *FieldLinks {
	return NewFieldLinks(
		contentType.client,
		fmt.Sprintf("%s/FieldLinks", contentType.endpoint),
		contentType.config,
	)
}
