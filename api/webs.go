package api

import (
	"encoding/json"
	"fmt"

	"github.com/koltyakov/gosip"
)

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

// Conf receives custom request config definition, e.g. custom headers, custom OData mod
func (webs *Webs) Conf(config *RequestConfig) *Webs {
	webs.config = config
	return webs
}

// Select adds $select OData modifier
func (webs *Webs) Select(oDataSelect string) *Webs {
	webs.modifiers.AddSelect(oDataSelect)
	return webs
}

// Expand adds $expand OData modifier
func (webs *Webs) Expand(oDataExpand string) *Webs {
	webs.modifiers.AddExpand(oDataExpand)
	return webs
}

// Filter adds $filter OData modifier
func (webs *Webs) Filter(oDataFilter string) *Webs {
	webs.modifiers.AddFilter(oDataFilter)
	return webs
}

// Top adds $top OData modifier
func (webs *Webs) Top(oDataTop int) *Webs {
	webs.modifiers.AddTop(oDataTop)
	return webs
}

// OrderBy adds $orderby OData modifier
func (webs *Webs) OrderBy(oDataOrderBy string, ascending bool) *Webs {
	webs.modifiers.AddOrderBy(oDataOrderBy, ascending)
	return webs
}

// Get ...
func (webs *Webs) Get() (WebsResp, error) {
	sp := NewHTTPClient(webs.client)
	headers := map[string]string{}
	if webs.config != nil {
		headers = webs.config.Headers
	}
	return sp.Get(webs.ToURL(), headers)
}

// Add ...
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

	body := trimMultiline(`{
		"parameters": ` + fmt.Sprintf("%s", parameters) + `
	}`)

	sp := NewHTTPClient(webs.client)
	headers := getConfHeaders(webs.config)

	headers["Accept"] = "application/json;odata=verbose"
	headers["Content-Type"] = "application/json;odata=verbose;charset=utf-8"

	return sp.Post(endpoint, []byte(body), headers)
}

/* Response helpers */

// Data : to get typed data
func (websResp *WebsResp) Data() []WebResp {
	collection, _ := parseODataCollection(*websResp)
	webs := []WebResp{}
	for _, web := range collection {
		webs = append(webs, WebResp(web))
	}
	return webs
}

// Unmarshal : to unmarshal to custom object
func (websResp *WebsResp) Unmarshal(obj interface{}) error {
	// collection := parseODataCollection(*websResp)
	// data, _ := json.Marshal(collection)
	data, _ := parseODataCollectionPlain(*websResp)
	return json.Unmarshal(data, obj)
}
