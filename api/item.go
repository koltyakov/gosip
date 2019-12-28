package api

import (
	"encoding/json"
	"fmt"
	"net/url"
	"time"

	"github.com/koltyakov/gosip"
)

// Item ...
type Item struct {
	client    *gosip.SPClient
	config    *RequestConfig
	endpoint  string
	modifiers map[string]string
}

// GenericItemInfo ...
type GenericItemInfo struct {
	ID            int       `json:"Id"`
	Title         string    `json:"Title"`
	ContentTypeID string    `json:"ContentTypeId"`
	Attachments   bool      `json:"Attachments"`
	AuthorID      int       `json:"AuthorId"`
	EditorID      int       `json:"EditorId"`
	Created       time.Time `json:"Created"`
	Modified      time.Time `json:"Modified"`
}

// ItemResp ...
type ItemResp []byte

// NewItem ...
func NewItem(client *gosip.SPClient, endpoint string, config *RequestConfig) *Item {
	return &Item{
		client:   client,
		endpoint: endpoint,
		config:   config,
	}
}

// ToURL ...
func (item *Item) ToURL() string {
	apiURL, _ := url.Parse(item.endpoint)
	query := apiURL.Query() // url.Values{}
	for k, v := range item.modifiers {
		query.Set(k, trimMultiline(v))
	}
	apiURL.RawQuery = query.Encode()
	return apiURL.String()
}

// Conf ...
func (item *Item) Conf(config *RequestConfig) *Item {
	item.config = config
	return item
}

// Select ...
func (item *Item) Select(oDataSelect string) *Item {
	if item.modifiers == nil {
		item.modifiers = make(map[string]string)
	}
	item.modifiers["$select"] = oDataSelect
	return item
}

// Expand ...
func (item *Item) Expand(oDataExpand string) *Item {
	if item.modifiers == nil {
		item.modifiers = make(map[string]string)
	}
	item.modifiers["$expand"] = oDataExpand
	return item
}

// Get ...
func (item *Item) Get() (ItemResp, error) {
	sp := NewHTTPClient(item.client)
	return sp.Get(item.ToURL(), getConfHeaders(item.config))
}

// Delete ...
func (item *Item) Delete() ([]byte, error) {
	sp := NewHTTPClient(item.client)
	return sp.Delete(item.endpoint, getConfHeaders(item.config))
}

// Update ...
func (item *Item) Update(body []byte) ([]byte, error) {
	body = patchMetadataTypeCB(body, func() string {
		endpoint := getPriorEndpoint(item.endpoint, "/Items")
		list := NewList(item.client, endpoint, nil)
		oDataType, _ := list.GetEntityType()
		return oDataType
	})
	sp := NewHTTPClient(item.client)
	return sp.Update(item.endpoint, body, getConfHeaders(item.config))
}

// Recycle ...
func (item *Item) Recycle() ([]byte, error) {
	endpoint := fmt.Sprintf("%s/Recycle", item.endpoint)
	sp := NewHTTPClient(item.client)
	return sp.Post(endpoint, nil, getConfHeaders(item.config))
}

// Roles ...
func (item *Item) Roles() *Roles {
	return NewRoles(item.client, item.endpoint, item.config)
}

// Attachments ...
func (item *Item) Attachments() *Attachments {
	return NewAttachments(
		item.client,
		fmt.Sprintf("%s/AttachmentFiles", item.endpoint),
		item.config,
	)
}

// ParentList ...
func (item *Item) ParentList() *List {
	return NewList(
		item.client,
		fmt.Sprintf("%s/ParentList", item.endpoint),
		item.config,
	)
}

/* Response helpers */

// Data : to get typed data
func (itemResp *ItemResp) Data() *GenericItemInfo {
	data := parseODataItem(*itemResp)
	data = fixDatesInResponse(data, []string{"Created", "Modified"})
	res := &GenericItemInfo{}
	json.Unmarshal(data, &res)
	return res
}

// Unmarshal : to unmarshal to custom object
func (itemResp *ItemResp) Unmarshal(obj interface{}) error {
	data := parseODataItem(*itemResp)
	data = normalizeMultiLookups(data)
	return json.Unmarshal(data, obj)
}
