package api

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/koltyakov/gosip"
)

//go:generate ggen -ent Fields -item Field -conf -coll -mods Select,Expand,Filter,Top,Skip,OrderBy -helpers Data,Normalized

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

// Get gets fieds response collection
func (fields *Fields) Get() (FieldsResp, error) {
	client := NewHTTPClient(fields.client)
	return client.Get(fields.ToURL(), fields.config)
}

// Add adds field with properties provided in `body` parameter
// where `body` is byte array representation of JSON string payload relevalt to SP.Field object
func (fields *Fields) Add(body []byte) (FieldResp, error) {
	body = patchMetadataType(body, "SP.Field")
	client := NewHTTPClient(fields.client)
	return client.Post(fields.endpoint, bytes.NewBuffer(body), fields.config)
}

// CreateFieldAsXML creates a field using XML schema definition
// `options` parameter (https://github.com/pnp/pnpjs/blob/version-2/packages/sp/fields/types.ts#L553)
// is only relevant for adding fields in list instances
func (fields *Fields) CreateFieldAsXML(schemaXML string, options int) (FieldResp, error) {
	endpoint := fmt.Sprintf("%s/CreateFieldAsXml", fields.endpoint)
	info := map[string]map[string]interface{}{
		"parameters": {
			"__metadata": &map[string]string{
				"type": "SP.XmlSchemaFieldCreationInformation",
			},
			"SchemaXml": schemaXML,
		},
	}
	if fields.entity == "list" {
		info["parameters"]["Options"] = options
	}
	payload, err := json.Marshal(info)
	if err != nil {
		return nil, err
	}
	client := NewHTTPClient(fields.client)
	return client.Post(endpoint, bytes.NewBuffer(payload), fields.config)
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

/* Pagination helpers */

// FieldsPage - paged items
type FieldsPage struct {
	Items       FieldsResp
	HasNextPage func() bool
	GetNextPage func() (*FieldsPage, error)
}

// GetPaged gets Paged Items collection
func (fields *Fields) GetPaged() (*FieldsPage, error) {
	data, err := fields.Get()
	if err != nil {
		return nil, err
	}
	res := &FieldsPage{
		Items: data,
		HasNextPage: func() bool {
			return data.HasNextPage()
		},
		GetNextPage: func() (*FieldsPage, error) {
			nextURL := data.NextPageURL()
			if nextURL == "" {
				return nil, fmt.Errorf("unable to get next page")
			}
			return NewFields(fields.client, nextURL, fields.config, fields.entity).GetPaged()
		},
	}
	return res, nil
}

// NextPageURL gets next page OData collection
func (fieldsResp *FieldsResp) NextPageURL() string {
	return getODataCollectionNextPageURL(*fieldsResp)
}

// HasNextPage returns is true if next page exists
func (fieldsResp *FieldsResp) HasNextPage() bool {
	return fieldsResp.NextPageURL() != ""
}
