package api

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/koltyakov/gosip"
)

//go:generate ggen -ent View -conf -mods Select,Expand -helpers Data,Normalized

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
	client := NewHTTPClient(view.client)
	return client.Get(view.ToURL(), view.config)
}

// Update updates View's metadata with properties provided in `body` parameter
// where `body` is byte array representation of JSON string payload relevalt to SP.View object
func (view *View) Update(body []byte) (ViewResp, error) {
	body = patchMetadataType(body, "SP.View")
	client := NewHTTPClient(view.client)
	return client.Update(view.endpoint, bytes.NewBuffer(body), view.config)
}

// Delete deletes this View (can't be restored from a recycle bin)
func (view *View) Delete() error {
	client := NewHTTPClient(view.client)
	_, err := client.Delete(view.endpoint, view.config)
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
	client := NewHTTPClient(view.client)
	return client.Post(endpoint, bytes.NewBuffer(payload), view.config)
}
