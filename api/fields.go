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
	entity    string
}

// FieldsResp - fields response type with helper processor methods
type FieldsResp []byte

// NewFields - Fields struct constructor function
func NewFields(client *gosip.SPClient, endpoint string, config *RequestConfig, entity string) *Fields {
	return &Fields{
		client:    client,
		endpoint:  endpoint,
		config:    config,
		modifiers: NewODataMods(),
		entity:    entity,
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

// Select adds $select OData modifier
func (fields *Fields) Select(oDataSelect string) *Fields {
	fields.modifiers.AddSelect(oDataSelect)
	return fields
}

// Expand adds $expand OData modifier
func (fields *Fields) Expand(oDataExpand string) *Fields {
	fields.modifiers.AddExpand(oDataExpand)
	return fields
}

// Filter adds $filter OData modifier
func (fields *Fields) Filter(oDataFilter string) *Fields {
	fields.modifiers.AddFilter(oDataFilter)
	return fields
}

// Top adds $top OData modifier
func (fields *Fields) Top(oDataTop int) *Fields {
	fields.modifiers.AddTop(oDataTop)
	return fields
}

// OrderBy adds $orderby OData modifier
func (fields *Fields) OrderBy(oDataOrderBy string, ascending bool) *Fields {
	fields.modifiers.AddOrderBy(oDataOrderBy, ascending)
	return fields
}

// Get gets fieds response collection
func (fields *Fields) Get() (FieldsResp, error) {
	sp := NewHTTPClient(fields.client)
	return sp.Get(fields.ToURL(), getConfHeaders(fields.config))
}

// Add adds field with properties provided in `body` parameter
// where `body` is byte array representation of JSON string payload relevalt to SP.Field object
func (fields *Fields) Add(body []byte) (FieldResp, error) {
	body = patchMetadataType(body, "SP.Field")
	sp := NewHTTPClient(fields.client)
	return sp.Post(fields.endpoint, body, getConfHeaders(fields.config))
}

// CreateFieldAsXML creates a field using XML schema definition
// `options` parameter (https://github.com/pnp/pnpjs/blob/version-2/packages/sp/fields/types.ts#L553)
// is only relevant for adding fields in list instances
func (fields *Fields) CreateFieldAsXML(schemaXML string, options int) (FieldResp, error) {
	endpoint := fmt.Sprintf("%s/CreateFieldAsXml", fields.endpoint)
	info := map[string]interface{}{
		"__metadata": &map[string]string{
			"type": "SP.XmlSchemaFieldCreationInformation",
		},
		"SchemaXml": schemaXML,
	}
	if fields.entity == "list" {
		info["Options"] = options
	}
	payload, err := json.Marshal(info)
	if err != nil {
		return nil, err
	}
	sp := NewHTTPClient(fields.client)
	return sp.Post(endpoint, payload, getConfHeaders(fields.config))
}

// GetByID gets a field by its ID (GUID)
func (fields *Fields) GetByID(fieldID string) *Field {
	return NewField(
		fields.client,
		fmt.Sprintf("%s('%s')", fields.endpoint, fieldID),
		fields.config,
	)
}

// GetByTitle gets a field by its Display Name
func (fields *Fields) GetByTitle(title string) *Field {
	return NewField(
		fields.client,
		fmt.Sprintf("%s/GetByTitle('%s')", fields.endpoint, title),
		fields.config,
	)
}

// GetByInternalNameOrTitle gets a field by its Internal or Display name
func (fields *Fields) GetByInternalNameOrTitle(internalName string) *Field {
	return NewField(
		fields.client,
		fmt.Sprintf("%s/GetByInternalNameOrTitle('%s')", fields.endpoint, internalName),
		fields.config,
	)
}

/* Response helpers */

// Data : to get typed data
func (fieldsResp *FieldsResp) Data() []FieldResp {
	collection, _ := normalizeODataCollection(*fieldsResp)
	resFields := []FieldResp{}
	for _, f := range collection {
		resFields = append(resFields, FieldResp(f))
	}
	return resFields
}

// Normalized returns normalized body
func (fieldsResp *FieldsResp) Normalized() []byte {
	normalized, _ := NormalizeODataCollection(*fieldsResp)
	return normalized
}
