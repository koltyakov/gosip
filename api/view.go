package api

import (
	"encoding/json"
	"fmt"

	"github.com/koltyakov/gosip"
)

// View represents SharePoint List View API queryable object struct
// Always use NewView constructor instead of &View{}
type View struct {
	client    *gosip.SPClient
	config    *RequestConfig
	endpoint  string
	modifiers *ODataMods
}

// ViewInfo - list view API response payload structure
type ViewInfo struct {
	BaseViewID                string `json:"BaseViewId"`
	DefaultView               bool   `json:"DefaultView"`
	DefaultViewForContentType bool   `json:"DefaultViewForContentType"`
	EditorModified            bool   `json:"EditorModified"`
	Hidden                    bool   `json:"Hidden"`
	HTMLSchemaXML             string `json:"HtmlSchemaXml"`
	ID                        string `json:"Id"`
	ImageURL                  string `json:"ImageUrl"`
	IncludeRootFolder         bool   `json:"IncludeRootFolder"`
	JSLink                    string `json:"JSLink"`
	ListViewXML               string `json:"ListViewXml"`
	MobileDefaultView         bool   `json:"MobileDefaultView"`
	MobileView                bool   `json:"MobileView"`
	OrderedView               bool   `json:"OrderedView"`
	Paged                     bool   `json:"Paged"`
	PersonalView              bool   `json:"PersonalView"`
	ReadOnlyView              bool   `json:"ReadOnlyView"`
	RequiresClientIntegration bool   `json:"RequiresClientIntegration"`
	RowLimit                  int    `json:"RowLimit"`
	Scope                     int    `json:"Scope"`
	ServerRelativeURL         string `json:"ServerRelativeUrl"`
	TabularView               bool   `json:"TabularView"`
	Threaded                  bool   `json:"Threaded"`
	Title                     string `json:"Title"`
	ViewQuery                 string `json:"ViewQuery"`
	ViewType                  string `json:"ViewType"`
}

// ViewResp - list view response type with helper processor methods
type ViewResp []byte

// NewView - View struct constructor function
func NewView(client *gosip.SPClient, endpoint string, config *RequestConfig) *View {
	return &View{
		client:    client,
		endpoint:  endpoint,
		config:    config,
		modifiers: NewODataMods(),
	}
}

// ToURL gets endpoint with modificators raw URL
func (view *View) ToURL() string {
	return toURL(view.endpoint, view.modifiers)
}

// Conf receives custom request config definition, e.g. custom headers, custom OData mod
func (view *View) Conf(config *RequestConfig) *View {
	view.config = config
	return view
}

// Select ...
func (view *View) Select(oDataSelect string) *View {
	view.modifiers.AddSelect(oDataSelect)
	return view
}

// Expand ...
func (view *View) Expand(oDataExpand string) *View {
	view.modifiers.AddExpand(oDataExpand)
	return view
}

// Get ...
func (view *View) Get() (ViewResp, error) {
	sp := NewHTTPClient(view.client)
	return sp.Get(view.ToURL(), getConfHeaders(view.config))
}

// Delete ...
func (view *View) Delete() ([]byte, error) {
	sp := NewHTTPClient(view.client)
	return sp.Delete(view.endpoint, getConfHeaders(view.config))
}

// Recycle ...
func (view *View) Recycle() ([]byte, error) {
	sp := NewHTTPClient(view.client)
	endpoint := fmt.Sprintf("%s/Recycle", view.endpoint)
	return sp.Post(endpoint, nil, getConfHeaders(view.config))
}

/* Response helpers */

// Data : to get typed data
func (viewResp *ViewResp) Data() *ViewInfo {
	data := parseODataItem(*viewResp)
	res := &ViewInfo{}
	json.Unmarshal(data, &res)
	return res
}

// Unmarshal : to unmarshal to custom object
func (viewResp *ViewResp) Unmarshal(obj interface{}) error {
	data := parseODataItem(*viewResp)
	return json.Unmarshal(data, obj)
}
