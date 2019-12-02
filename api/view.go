package api

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/koltyakov/gosip"
)

// View ...
type View struct {
	client    *gosip.SPClient
	config    *RequestConfig
	endpoint  string
	modifiers map[string]string
}

// ViewInfo ...
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

// ViewResp ...
type ViewResp []byte

// NewView ...
func NewView(client *gosip.SPClient, endpoint string, config *RequestConfig) *View {
	return &View{
		client:   client,
		endpoint: endpoint,
		config:   config,
	}
}

// ToURL ...
func (view *View) ToURL() string {
	apiURL, _ := url.Parse(view.endpoint)
	query := url.Values{}
	for k, v := range view.modifiers {
		query.Add(k, trimMultiline(v))
	}
	apiURL.RawQuery = query.Encode()
	return apiURL.String()
}

// Conf ...
func (view *View) Conf(config *RequestConfig) *View {
	view.config = config
	return view
}

// Select ...
func (view *View) Select(oDataSelect string) *View {
	if view.modifiers == nil {
		view.modifiers = make(map[string]string)
	}
	view.modifiers["$select"] = oDataSelect
	return view
}

// Expand ...
func (view *View) Expand(oDataExpand string) *View {
	if view.modifiers == nil {
		view.modifiers = make(map[string]string)
	}
	view.modifiers["$expand"] = oDataExpand
	return view
}

// Get ...
func (view *View) Get() (ViewResp, error) {
	sp := NewHTTPClient(view.client)
	data, err := sp.Get(view.ToURL(), getConfHeaders(view.config))
	if err != nil {
		return nil, err
	}
	return ViewResp(data), nil
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
func (viewResp *ViewResp) Unmarshal(obj *interface{}) error {
	data := parseODataItem(*viewResp)
	return json.Unmarshal(data, &obj)
}
