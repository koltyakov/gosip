package api

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/koltyakov/gosip"
)

// Fields represent SharePoint Fields (Site Columns) API queryable collection struct
type Fields struct {
	client    *gosip.SPClient
	config    *RequestConfig
	endpoint  string
	modifiers map[string]string
}

// FieldsResp - fields response type with helper processor methods
type FieldsResp []byte

// NewFields - Fields struct constructor function
func NewFields(client *gosip.SPClient, endpoint string, config *RequestConfig) *Fields {
	return &Fields{
		client:   client,
		endpoint: endpoint,
		config:   config,
	}
}

// ToURL ...
func (fields *Fields) ToURL() string {
	apiURL, _ := url.Parse(fields.endpoint)
	query := apiURL.Query() // url.Values{}
	for k, v := range fields.modifiers {
		query.Set(k, trimMultiline(v))
	}
	apiURL.RawQuery = query.Encode()
	return apiURL.String()
}

// Conf ...
func (fields *Fields) Conf(config *RequestConfig) *Fields {
	fields.config = config
	return fields
}

// Select ...
func (fields *Fields) Select(oDataSelect string) *Fields {
	if fields.modifiers == nil {
		fields.modifiers = make(map[string]string)
	}
	fields.modifiers["$select"] = oDataSelect
	return fields
}

// Expand ...
func (fields *Fields) Expand(oDataExpand string) *Fields {
	if fields.modifiers == nil {
		fields.modifiers = make(map[string]string)
	}
	fields.modifiers["$expand"] = oDataExpand
	return fields
}

// Filter ...
func (fields *Fields) Filter(oDataFilter string) *Fields {
	if fields.modifiers == nil {
		fields.modifiers = make(map[string]string)
	}
	fields.modifiers["$filter"] = oDataFilter
	return fields
}

// Top ...
func (fields *Fields) Top(oDataTop int) *Fields {
	if fields.modifiers == nil {
		fields.modifiers = make(map[string]string)
	}
	fields.modifiers["$top"] = fmt.Sprintf("%d", oDataTop)
	return fields
}

// OrderBy ...
func (fields *Fields) OrderBy(oDataOrderBy string, ascending bool) *Fields {
	direction := "asc"
	if !ascending {
		direction = "desc"
	}
	if fields.modifiers == nil {
		fields.modifiers = make(map[string]string)
	}
	if fields.modifiers["$orderby"] != "" {
		fields.modifiers["$orderby"] += ","
	}
	fields.modifiers["$orderby"] += fmt.Sprintf("%s %s", oDataOrderBy, direction)
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
