package api

import (
	"encoding/json"
	"fmt"

	"github.com/koltyakov/gosip"
)

// Field represents SharePoint Field (Site Column) API queryable object struct
// Always use NewField constructor instead of &Field{}
type Field struct {
	client    *gosip.SPClient
	config    *RequestConfig
	endpoint  string
	modifiers *ODataMods
}

// GenericFieldInfo - generic field API response payload structure
type GenericFieldInfo struct {
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

// Conf receives custom request config definition, e.g. custom headers, custom OData mod
func (field *Field) Conf(config *RequestConfig) *Field {
	field.config = config
	return field
}

// Select ...
func (field *Field) Select(oDataSelect string) *Field {
	field.modifiers.AddSelect(oDataSelect)
	return field
}

// Expand ...
func (field *Field) Expand(oDataExpand string) *Field {
	field.modifiers.AddExpand(oDataExpand)
	return field
}

// Get ...
func (field *Field) Get() (FieldResp, error) {
	sp := NewHTTPClient(field.client)
	data, err := sp.Get(field.ToURL(), getConfHeaders(field.config))
	if err != nil {
		return nil, err
	}
	return FieldResp(data), nil
}

// Delete ...
func (field *Field) Delete() ([]byte, error) {
	sp := NewHTTPClient(field.client)
	return sp.Delete(field.endpoint, getConfHeaders(field.config))
}

// Recycle ...
func (field *Field) Recycle() ([]byte, error) {
	sp := NewHTTPClient(field.client)
	endpoint := fmt.Sprintf("%s/Recycle", field.endpoint)
	return sp.Post(endpoint, nil, getConfHeaders(field.config))
}

/* Response helpers */

// Data : to get typed data
func (fieldResp *FieldResp) Data() *GenericFieldInfo {
	data := parseODataItem(*fieldResp)
	res := &GenericFieldInfo{}
	json.Unmarshal(data, &res)
	return res
}

// Unmarshal : to unmarshal to custom object
func (fieldResp *FieldResp) Unmarshal(obj interface{}) error {
	data := parseODataItem(*fieldResp)
	return json.Unmarshal(data, obj)
}
