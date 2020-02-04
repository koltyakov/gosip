package api

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/koltyakov/gosip"
)

//go:generate ggen -ent Webs -item Web -conf -coll -mods Select,Expand,Filter,Top,OrderBy -helpers Data,Normalized

// Webs represent SharePoint Webs API queryable collection struct
// Always use NewWebs constructor instead of &Webs{}
type Webs struct {
	client    *gosip.SPClient
	config    *RequestConfig
	endpoint  string
	modifiers *ODataMods
}

// WebsResp - webs response type with helper processor methods
type WebsResp []byte

// NewWebs - Webs struct constructor function
func NewWebs(client *gosip.SPClient, endpoint string, config *RequestConfig) *Webs {
	return &Webs{
		client:    client,
		endpoint:  endpoint,
		config:    config,
		modifiers: NewODataMods(),
	}
}

// ToURL gets endpoint with modificators raw URL
func (webs *Webs) ToURL() string {
	return toURL(webs.endpoint, webs.modifiers)
}

// Get gets Webs response - a collection of WebInfo for the parent Web
func (webs *Webs) Get() (WebsResp, error) {
	sp := NewHTTPClient(webs.client)
	headers := map[string]string{}
	if webs.config != nil {
		headers = webs.config.Headers
	}
	return sp.Get(webs.ToURL(), headers)
}

// Add creates a subweb for a parent web with provided `title` and `url`.
// `url` stands for a system friendly URI (e.g. `finances`) while `title` is a human friendly name (e.g. `Financial Department`).
// Along with title and url additional metadata can be provided in optional `metadata` string map object.
// `metadata` props should correspond to `SP.WebCreationInformation` API type. Some props have defaults as Language (1033), WebTemplate (STS), etc.
func (webs *Webs) Add(title string, url string, metadata map[string]interface{}) (WebResp, error) {
	endpoint := fmt.Sprintf("%s/Add", webs.endpoint)

	if metadata == nil {
		metadata = make(map[string]interface{})
	}

	metadata["__metadata"] = map[string]string{
		"type": "SP.WebCreationInformation",
	}

	metadata["Title"] = title
	metadata["Url"] = url

	// metadata["Description"]

	// Default values
	if metadata["Language"] == nil {
		metadata["Language"] = 1033
	}
	if metadata["UseSamePermissionsAsParentSite"] == nil {
		metadata["UseSamePermissionsAsParentSite"] = true
	}
	if metadata["WebTemplate"] == nil {
		metadata["WebTemplate"] = "STS"
	}

	parameters, _ := json.Marshal(metadata)

	body := TrimMultiline(`{
		"parameters": ` + fmt.Sprintf("%s", parameters) + `
	}`)

	sp := NewHTTPClient(webs.client)
	headers := getConfHeaders(webs.config)

	headers["Accept"] = "application/json;odata=verbose"
	headers["Content-Type"] = "application/json;odata=verbose;charset=utf-8"

	return sp.Post(endpoint, bytes.NewBuffer([]byte(body)), headers)
}
