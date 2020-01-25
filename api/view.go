package api

import (
	"encoding/json"
	"fmt"

	"github.com/koltyakov/gosip"
)

//go:generate ggen -ent View -conf -mods Select,Expand

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

// Get gets this View data response
func (view *View) Get() (ViewResp, error) {
	sp := NewHTTPClient(view.client)
	return sp.Get(view.ToURL(), getConfHeaders(view.config))
}

// Update updates View's metadata with properties provided in `body` parameter
// where `body` is byte array representation of JSON string payload relevalt to SP.View object
func (view *View) Update(body []byte) (ViewResp, error) {
	body = patchMetadataType(body, "SP.View")
	sp := NewHTTPClient(view.client)
	return sp.Update(view.endpoint, body, getConfHeaders(view.config))
}

// Delete deletes this View (can't be restored from a recycle bin)
func (view *View) Delete() error {
	sp := NewHTTPClient(view.client)
	_, err := sp.Delete(view.endpoint, getConfHeaders(view.config))
	return err
}

// SetViewXML updates view XML
func (view *View) SetViewXML(viewXML string) (ViewResp, error) {
	endpoint := fmt.Sprintf("%s/SetViewXml()", view.endpoint)
	payload, err := json.Marshal(&struct {
		ViewXML string `json:"viewXml"`
	}{
		ViewXML: viewXML,
	})
	if err != nil {
		return nil, err
	}
	sp := NewHTTPClient(view.client)
	return sp.Post(endpoint, payload, getConfHeaders(view.config))
}

/* Response helpers */

// Data : to get typed data
func (viewResp *ViewResp) Data() *ViewInfo {
	data := NormalizeODataItem(*viewResp)
	res := &ViewInfo{}
	json.Unmarshal(data, &res)
	return res
}

// Normalized returns normalized body
func (viewResp *ViewResp) Normalized() []byte {
	return NormalizeODataItem(*viewResp)
}
