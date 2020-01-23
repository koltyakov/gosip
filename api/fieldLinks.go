package api

import (
	"encoding/json"

	"github.com/koltyakov/gosip"
)

// FieldLinks represent SharePoint content type FieldLinks API queryable collection struct
// Always use NewFieldLinks constructor instead of &FieldLinks{}
type FieldLinks struct {
	client    *gosip.SPClient
	config    *RequestConfig
	endpoint  string
	modifiers *ODataMods
}

// FieldLinkInfo field link info
type FieldLinkInfo struct {
	ID                string `json:"Id"`
	Name              string `json:"Name"`
	FieldInternalName string `json:"FieldInternalName"`
	Hidden            bool   `json:"Hidden"`
	Required          bool   `json:"Required"`
}

// FieldLinksResp - fieldLinks response type with helper processor methods
type FieldLinksResp []byte

// FieldLinkResp - fieldLinks response type with helper processor methods
type FieldLinkResp []byte

// NewFieldLinks - FieldLinks struct constructor function
func NewFieldLinks(client *gosip.SPClient, endpoint string, config *RequestConfig) *FieldLinks {
	return &FieldLinks{
		client:    client,
		endpoint:  endpoint,
		config:    config,
		modifiers: NewODataMods(),
	}
}

// ToURL gets endpoint with modificators raw URL
func (fieldLinks *FieldLinks) ToURL() string {
	return toURL(fieldLinks.endpoint, fieldLinks.modifiers)
}

// Conf receives custom request config definition, e.g. custom headers, custom OData mod
func (fieldLinks *FieldLinks) Conf(config *RequestConfig) *FieldLinks {
	fieldLinks.config = config
	return fieldLinks
}

// Filter adds $filter OData modifier
func (fieldLinks *FieldLinks) Filter(oDataFilter string) *FieldLinks {
	fieldLinks.modifiers.AddFilter(oDataFilter)
	return fieldLinks
}

// Top adds $top OData modifier
func (fieldLinks *FieldLinks) Top(oDataTop int) *FieldLinks {
	fieldLinks.modifiers.AddTop(oDataTop)
	return fieldLinks
}

// Get gets fieds response collection
func (fieldLinks *FieldLinks) Get() (FieldLinksResp, error) {
	sp := NewHTTPClient(fieldLinks.client)
	return sp.Get(fieldLinks.ToURL(), getConfHeaders(fieldLinks.config))
}

// // Add adds field link
// func (fieldLinks *FieldLinks) Add(name string, hidden bool, required bool) (*FieldLinkInfo, error) {
// 	// REST API doesn't work in that context as supposed to (https://social.msdn.microsoft.com/Forums/office/en-US/52dc9d24-2eb3-4540-a26a-02b12fe1155b/rest-add-fieldlink-to-content-type?forum=sharepointdevelopment)
// 	body := []byte(TrimMultiline(fmt.Sprintf(`{
// 		"__metadata": { "type": "SP.FieldLink" },
// 		"FieldInternalName": "%s",
// 		"Hidden": %t,
// 		"Required": %t
// 	}`, name, hidden, required)))
// 	sp := NewHTTPClient(fieldLinks.client)
// 	resp, err := sp.Post(fieldLinks.endpoint, body, getConfHeaders(fieldLinks.config))
// 	if err != nil {
// 		return nil, err
// 	}
// 	linkInfo := &FieldLinkInfo{}
// 	if err := json.Unmarshal(resp, &linkInfo); err != nil {
// 		return nil, err
// 	}
// 	return linkInfo, nil
// }

/* Response helpers */

// Data : to get typed data
func (fieldLinksResp *FieldLinksResp) Data() []*FieldLinkInfo {
	collection, _ := normalizeODataCollection(*fieldLinksResp)
	resFieldLinks := []*FieldLinkInfo{}
	for _, f := range collection {
		linkInfo := &FieldLinkInfo{}
		json.Unmarshal(f, &linkInfo)
		resFieldLinks = append(resFieldLinks, linkInfo)
	}
	return resFieldLinks
}

// Normalized returns normalized body
func (fieldLinksResp *FieldLinksResp) Normalized() []byte {
	normalized, _ := NormalizeODataCollection(*fieldLinksResp)
	return normalized
}
