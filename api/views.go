package api

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/koltyakov/gosip"
)

// Views ...
type Views struct {
	client    *gosip.SPClient
	config    *RequestConfig
	endpoint  string
	modifiers map[string]string
}

// ViewsResp ...
type ViewsResp []byte

// NewViews ...
func NewViews(client *gosip.SPClient, endpoint string, config *RequestConfig) *Views {
	return &Views{
		client:   client,
		endpoint: endpoint,
		config:   config,
	}
}

// ToURL ...
func (views *Views) ToURL() string {
	apiURL, _ := url.Parse(views.endpoint)
	query := url.Values{}
	for k, v := range views.modifiers {
		query.Add(k, trimMultiline(v))
	}
	apiURL.RawQuery = query.Encode()
	return apiURL.String()
}

// Conf ...
func (views *Views) Conf(config *RequestConfig) *Views {
	views.config = config
	return views
}

// Select ...
func (views *Views) Select(oDataSelect string) *Views {
	if views.modifiers == nil {
		views.modifiers = make(map[string]string)
	}
	views.modifiers["$select"] = oDataSelect
	return views
}

// Expand ...
func (views *Views) Expand(oDataExpand string) *Views {
	if views.modifiers == nil {
		views.modifiers = make(map[string]string)
	}
	views.modifiers["$expand"] = oDataExpand
	return views
}

// Filter ...
func (views *Views) Filter(oDataFilter string) *Views {
	if views.modifiers == nil {
		views.modifiers = make(map[string]string)
	}
	views.modifiers["$filter"] = oDataFilter
	return views
}

// Top ...
func (views *Views) Top(oDataTop int) *Views {
	if views.modifiers == nil {
		views.modifiers = make(map[string]string)
	}
	views.modifiers["$top"] = fmt.Sprintf("%d", oDataTop)
	return views
}

// OrderBy ...
func (views *Views) OrderBy(oDataOrderBy string, ascending bool) *Views {
	direction := "asc"
	if !ascending {
		direction = "desc"
	}
	if views.modifiers == nil {
		views.modifiers = make(map[string]string)
	}
	if views.modifiers["$orderby"] != "" {
		views.modifiers["$orderby"] += ","
	}
	views.modifiers["$orderby"] += fmt.Sprintf("%s %s", oDataOrderBy, direction)
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
	collection := parseODataCollection(*viewsResp)
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
	data := parseODataCollectionPlain(*viewsResp)
	return json.Unmarshal(data, obj)
}
