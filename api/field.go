package api

import (
	"bytes"

	"github.com/koltyakov/gosip"
)

//go:generate ggen -ent Field -conf -mods Select,Expand -helpers Data,Normalized

// Field represents SharePoint Field (Site Column) API queryable object struct
// Always use NewField constructor instead of &Field{}
type Field struct {
	client    *gosip.SPClient
	config    *RequestConfig
	endpoint  string
	modifiers *ODataMods
}

// FieldInfo - generic field API response payload structure
type FieldInfo struct {
	AutoIndexed          bool   `json:"AutoIndexed"`
	CanBeDeleted         bool   `json:"CanBeDeleted"`
	DefaultValue         string `json:"DefaultValue"`
	Description          string `json:"Description"`
	EnforceUniqueValues  bool   `json:"EnforceUniqueValues"`
	EntityPropertyName   string `json:"EntityPropertyName"`
	FieldTypeKind        int    `json:"FieldTypeKind"`
	Filterable           bool   `json:"Filterable"`
	FromBaseType         bool   `json:"FromBaseType"`
	Group                string `json:"Group"`
	Hidden               bool   `json:"Hidden"`
	ID                   string `json:"Id"`
	IndexStatus          int    `json:"IndexStatus"`
	Indexed              bool   `json:"Indexed"`
	InternalName         string `json:"InternalName"`
	JSLink               string `json:"JSLink"`
	LookupField          string `json:"LookupField"`
	LookupList           string `json:"LookupList"`
	LookupWebID          string `json:"LookupWebId"`
	ReadOnlyField        bool   `json:"ReadOnlyField"`
	Required             bool   `json:"Required"`
	SchemaXML            string `json:"SchemaXml"`
	Scope                string `json:"Scope"`
	Sealed               bool   `json:"Sealed"`
	ShowInFiltersPane    int    `json:"ShowInFiltersPane"`
	Sortable             bool   `json:"Sortable"`
	StaticName           string `json:"StaticName"`
	Title                string `json:"Title"`
	TypeAsString         string `json:"TypeAsString"`
	TypeDisplayName      string `json:"TypeDisplayName"`
	TypeShortDescription string `json:"TypeShortDescription"`
}

// FieldResp - field response type with helper processor methods
type FieldResp []byte

// NewField - Field struct constructor function
func NewField(client *gosip.SPClient, endpoint string, config *RequestConfig) *Field {
	return &Field{
		client:    client,
		endpoint:  endpoint,
		config:    config,
		modifiers: NewODataMods(),
	}
}

// ToURL gets endpoint with modificators raw URL
func (field *Field) ToURL() string {
	return toURL(field.endpoint, field.modifiers)
}

// Get gets field data object
func (field *Field) Get() (FieldResp, error) {
	client := NewHTTPClient(field.client)
	data, err := client.Get(field.ToURL(), field.config)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// Update updates Field's metadata with properties provided in `body` parameter
// where `body` is byte array representation of JSON string payload relevant to SP.Field object
func (field *Field) Update(body []byte) (FieldResp, error) {
	body = patchMetadataType(body, "SP.Field")
	client := NewHTTPClient(field.client)
	return client.Update(field.endpoint, bytes.NewBuffer(body), field.config)
}

// Delete deletes a field skipping recycle bin
func (field *Field) Delete() error {
	client := NewHTTPClient(field.client)
	_, err := client.Delete(field.endpoint, field.config)
	return err
}
