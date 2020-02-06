package api

import (
	"bytes"
	"fmt"

	"github.com/koltyakov/gosip"
)

//go:generate ggen -ent Views -item View -conf -coll -mods Select,Expand,Filter,Top,Skip,OrderBy -helpers Data,Normalized,Pagination

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

// Get gets this List or Document Library views collection
func (views *Views) Get() (ViewsResp, error) {
	client := NewHTTPClient(views.client)
	return client.Get(views.ToURL(), views.config)
}

// Add adds view with properties provided in `body` parameter
// where `body` is byte array representation of JSON string payload relevalt to SP.View object
func (views *Views) Add(body []byte) (ViewResp, error) {
	body = patchMetadataType(body, "SP.View")
	client := NewHTTPClient(views.client)
	return client.Post(views.endpoint, bytes.NewBuffer(body), views.config)
}

// GetByID gets a view by its ID (GUID)
func (views *Views) GetByID(viewID string) *View {
	return NewView(
		views.client,
		fmt.Sprintf("%s('%s')", views.endpoint, viewID),
		views.config,
	)
}

// GetByTitle gets a view by its Display Name (Title)
func (views *Views) GetByTitle(title string) *View {
	return NewView(
		views.client,
		fmt.Sprintf("%s/GetByTitle('%s')", views.endpoint, title),
		views.config,
	)
}

// DefaultView gets list's default view
func (views *Views) DefaultView() *View {
	return NewView(
		views.client,
		fmt.Sprintf("%s/DefaultView", getPriorEndpoint(views.endpoint, "/Views")),
		views.config,
	)
}
