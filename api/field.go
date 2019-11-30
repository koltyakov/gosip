package api

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/koltyakov/gosip"
)

// Field ...
type Field struct {
	client    *gosip.SPClient
	config    *RequestConfig
	endpoint  string
	modifiers map[string]string
}

// FieldInfo ...
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
	ReadOnlyField        bool   `json:"ReadOnlyField"`
	Required             bool   `json:"Required"`
	SchemaXml            string `json:"SchemaXml"`
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

// NewField ...
func NewField(client *gosip.SPClient, endpoint string, config *RequestConfig) *Field {
	return &Field{
		client:   client,
		endpoint: endpoint,
		config:   config,
	}
}

// ToURL ...
func (field *Field) ToURL() string {
	apiURL, _ := url.Parse(field.endpoint)
	query := url.Values{}
	for k, v := range field.modifiers {
		query.Add(k, trimMultiline(v))
	}
	apiURL.RawQuery = query.Encode()
	return apiURL.String()
}

// Conf ...
func (field *Field) Conf(config *RequestConfig) *Field {
	field.config = config
	return field
}

// Select ...
func (field *Field) Select(oDataSelect string) *Field {
	if field.modifiers == nil {
		field.modifiers = make(map[string]string)
	}
	field.modifiers["$select"] = oDataSelect
	return field
}

// Expand ...
func (field *Field) Expand(oDataExpand string) *Field {
	if field.modifiers == nil {
		field.modifiers = make(map[string]string)
	}
	field.modifiers["$expand"] = oDataExpand
	return field
}

// Get ...
func (field *Field) Get() (*FieldInfo, error) {
	sp := NewHTTPClient(field.client)
	data, err := sp.Get(field.ToURL(), HeadersPresets.Verbose.Headers)
	if err != nil {
		return nil, err
	}
	res := &struct {
		Field *FieldInfo `json:"d"`
	}{}
	if err := json.Unmarshal(data, &res); err != nil {
		return nil, err
	}
	return res.Field, nil
}

// GetRaw ...
func (field *Field) GetRaw() ([]byte, error) {
	sp := NewHTTPClient(field.client)
	return sp.Get(field.ToURL(), getConfHeaders(field.config))
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
