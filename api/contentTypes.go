package api

import (
	"encoding/json"
	"fmt"

	"github.com/koltyakov/gosip"
)

// ContentTypes represent SharePoint Content Types API queryable collection struct
type ContentTypes struct {
	client    *gosip.SPClient
	config    *RequestConfig
	endpoint  string
	modifiers map[string]string
}

// ContentTypesResp - content types response type with helper processor methods
type ContentTypesResp []byte

// NewContentTypes - ContentTypes struct constructor function
func NewContentTypes(client *gosip.SPClient, endpoint string, config *RequestConfig) *ContentTypes {
	return &ContentTypes{
		client:   client,
		endpoint: endpoint,
		config:   config,
	}
}

// ToURL gets endpoint with modificators raw URL ...
func (contentTypes *ContentTypes) ToURL() string {
	return toURL(contentTypes.endpoint, contentTypes.modifiers)
}

// Conf ...
func (contentTypes *ContentTypes) Conf(config *RequestConfig) *ContentTypes {
	contentTypes.config = config
	return contentTypes
}

// Select ...
func (contentTypes *ContentTypes) Select(oDataSelect string) *ContentTypes {
	if contentTypes.modifiers == nil {
		contentTypes.modifiers = make(map[string]string)
	}
	contentTypes.modifiers["$select"] = oDataSelect
	return contentTypes
}

// Expand ...
func (contentTypes *ContentTypes) Expand(oDataExpand string) *ContentTypes {
	if contentTypes.modifiers == nil {
		contentTypes.modifiers = make(map[string]string)
	}
	contentTypes.modifiers["$expand"] = oDataExpand
	return contentTypes
}

// Filter ...
func (contentTypes *ContentTypes) Filter(oDataFilter string) *ContentTypes {
	if contentTypes.modifiers == nil {
		contentTypes.modifiers = make(map[string]string)
	}
	contentTypes.modifiers["$filter"] = oDataFilter
	return contentTypes
}

// Top ...
func (contentTypes *ContentTypes) Top(oDataTop int) *ContentTypes {
	if contentTypes.modifiers == nil {
		contentTypes.modifiers = make(map[string]string)
	}
	contentTypes.modifiers["$top"] = fmt.Sprintf("%d", oDataTop)
	return contentTypes
}

// OrderBy ...
func (contentTypes *ContentTypes) OrderBy(oDataOrderBy string, ascending bool) *ContentTypes {
	direction := "asc"
	if !ascending {
		direction = "desc"
	}
	if contentTypes.modifiers == nil {
		contentTypes.modifiers = make(map[string]string)
	}
	if contentTypes.modifiers["$orderby"] != "" {
		contentTypes.modifiers["$orderby"] += ","
	}
	contentTypes.modifiers["$orderby"] += fmt.Sprintf("%s %s", oDataOrderBy, direction)
	return contentTypes
}

// Get ...
func (contentTypes *ContentTypes) Get() (ContentTypesResp, error) {
	sp := NewHTTPClient(contentTypes.client)
	return sp.Get(contentTypes.ToURL(), getConfHeaders(contentTypes.config))
}

// GetByID ...
func (contentTypes *ContentTypes) GetByID(contentTypeID string) *ContentType {
	return NewContentType(
		contentTypes.client,
		fmt.Sprintf("%s('%s')", contentTypes.endpoint, contentTypeID),
		contentTypes.config,
	)
}

// ToDo:
// Add

/* Response helpers */

// Data : to get typed data
func (ctsResp *ContentTypesResp) Data() []ContentTypeResp {
	collection, _ := parseODataCollection(*ctsResp)
	cts := []ContentTypeResp{}
	for _, ct := range collection {
		cts = append(cts, ContentTypeResp(ct))
	}
	return cts
}

// Unmarshal : to unmarshal to custom object
func (ctsResp *ContentTypesResp) Unmarshal(obj interface{}) error {
	// collection := parseODataCollection(*ctsResp)
	// data, _ := json.Marshal(collection)
	data, _ := parseODataCollectionPlain(*ctsResp)
	return json.Unmarshal(data, obj)
}
