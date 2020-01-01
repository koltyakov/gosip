package api

import (
	"encoding/json"
	"fmt"

	"github.com/koltyakov/gosip"
)

// Views  represent SharePoint List Views API queryable collection struct
// Always use NewViews constructor instead of &Views{}
type Views struct {
	client    *gosip.SPClient
	config    *RequestConfig
	endpoint  string
	modifiers *ODataMods
}

// ViewsResp - list views response type with helper processor methods
type ViewsResp []byte

// NewViews - Views struct constructor function
func NewViews(client *gosip.SPClient, endpoint string, config *RequestConfig) *Views {
	return &Views{
		client:    client,
		endpoint:  endpoint,
		config:    config,
		modifiers: NewODataMods(),
	}
}

// ToURL gets endpoint with modificators raw URL
func (views *Views) ToURL() string {
	return toURL(views.endpoint, views.modifiers)
}

// Conf receives custom request config definition, e.g. custom headers, custom OData mod
func (views *Views) Conf(config *RequestConfig) *Views {
	views.config = config
	return views
}

// Select ...
func (views *Views) Select(oDataSelect string) *Views {
	views.modifiers.AddSelect(oDataSelect)
	return views
}

// Expand ...
func (views *Views) Expand(oDataExpand string) *Views {
	views.modifiers.AddExpand(oDataExpand)
	return views
}

// Filter ...
func (views *Views) Filter(oDataFilter string) *Views {
	views.modifiers.AddFilter(oDataFilter)
	return views
}

// Top ...
func (views *Views) Top(oDataTop int) *Views {
	views.modifiers.AddTop(oDataTop)
	return views
}

// OrderBy ...
func (views *Views) OrderBy(oDataOrderBy string, ascending bool) *Views {
	views.modifiers.AddOrderBy(oDataOrderBy, ascending)
	return views
}

// Get ...
func (views *Views) Get() (ViewsResp, error) {
	sp := NewHTTPClient(views.client)
	return sp.Get(views.ToURL(), getConfHeaders(views.config))
}

// GetByID ...
func (views *Views) GetByID(viewID string) *View {
	return NewView(
		views.client,
		fmt.Sprintf("%s('%s')", views.endpoint, viewID),
		views.config,
	)
}

// GetByTitle ...
func (views *Views) GetByTitle(title string) *View {
	return NewView(
		views.client,
		fmt.Sprintf("%s/GetByTitle('%s')", views.endpoint, title),
		views.config,
	)
}

// ToDo:
// Add

/* Response helpers */

// Data : to get typed data
func (viewsResp *ViewsResp) Data() []ViewResp {
	collection, _ := parseODataCollection(*viewsResp)
	views := []ViewResp{}
	for _, view := range collection {
		views = append(views, ViewResp(view))
	}
	return views
}

// Unmarshal : to unmarshal to custom object
func (viewsResp *ViewsResp) Unmarshal(obj interface{}) error {
	// collection := parseODataCollection(*viewsResp)
	// data, _ := json.Marshal(collection)
	data, _ := parseODataCollectionPlain(*viewsResp)
	return json.Unmarshal(data, obj)
}
