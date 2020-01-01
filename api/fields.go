package api

import (
	"encoding/json"
	"fmt"

	"github.com/koltyakov/gosip"
)

// Fields represent SharePoint Fields (Site Columns) API queryable collection struct
// Always use NewFields constructor instead of &Fields{}
type Fields struct {
	client    *gosip.SPClient
	config    *RequestConfig
	endpoint  string
	modifiers *ODataMods
}

// FieldsResp - fields response type with helper processor methods
type FieldsResp []byte

// NewFields - Fields struct constructor function
func NewFields(client *gosip.SPClient, endpoint string, config *RequestConfig) *Fields {
	return &Fields{
		client:    client,
		endpoint:  endpoint,
		config:    config,
		modifiers: NewODataMods(),
	}
}

// ToURL gets endpoint with modificators raw URL
func (fields *Fields) ToURL() string {
	return toURL(fields.endpoint, fields.modifiers)
}

// Conf receives custom request config definition, e.g. custom headers, custom OData mod
func (fields *Fields) Conf(config *RequestConfig) *Fields {
	fields.config = config
	return fields
}

// Select ...
func (fields *Fields) Select(oDataSelect string) *Fields {
	fields.modifiers.AddSelect(oDataSelect)
	return fields
}

// Expand ...
func (fields *Fields) Expand(oDataExpand string) *Fields {
	fields.modifiers.AddExpand(oDataExpand)
	return fields
}

// Filter ...
func (fields *Fields) Filter(oDataFilter string) *Fields {
	fields.modifiers.AddFilter(oDataFilter)
	return fields
}

// Top ...
func (fields *Fields) Top(oDataTop int) *Fields {
	fields.modifiers.AddTop(oDataTop)
	return fields
}

// OrderBy ...
func (fields *Fields) OrderBy(oDataOrderBy string, ascending bool) *Fields {
	fields.modifiers.AddOrderBy(oDataOrderBy, ascending)
	return fields
}

// Get ...
func (fields *Fields) Get() (FieldsResp, error) {
	sp := NewHTTPClient(fields.client)
	return sp.Get(fields.ToURL(), getConfHeaders(fields.config))
}

// GetByID ...
func (fields *Fields) GetByID(fieldID string) *Field {
	return NewField(
		fields.client,
		fmt.Sprintf("%s('%s')", fields.endpoint, fieldID),
		fields.config,
	)
}

// GetByTitle ...
func (fields *Fields) GetByTitle(title string) *Field {
	return NewField(
		fields.client,
		fmt.Sprintf("%s/GetByTitle('%s')", fields.endpoint, title),
		fields.config,
	)
}

// GetByInternalNameOrTitle ...
func (fields *Fields) GetByInternalNameOrTitle(internalName string) *Field {
	return NewField(
		fields.client,
		fmt.Sprintf("%s/GetByInternalNameOrTitle('%s')", fields.endpoint, internalName),
		fields.config,
	)
}

// ToDo:
// Add

/* Response helpers */

// Data : to get typed data
func (fieldsResp *FieldsResp) Data() []FieldResp {
	collection, _ := parseODataCollection(*fieldsResp)
	resFields := []FieldResp{}
	for _, f := range collection {
		resFields = append(resFields, FieldResp(f))
	}
	return resFields
}

// Unmarshal : to unmarshal to custom object
func (fieldsResp *FieldsResp) Unmarshal(obj interface{}) error {
	// collection := parseODataCollection(*fieldsResp)
	// data, _ := json.Marshal(collection)
	data, _ := parseODataCollectionPlain(*fieldsResp)
	return json.Unmarshal(data, obj)
}
