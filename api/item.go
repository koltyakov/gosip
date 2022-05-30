package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/koltyakov/gosip"
)

//go:generate ggen -ent Item -conf -mods Select,Expand -helpers Normalized,ToMap

// Item represents SharePoint Lists & Document Libraries Items API queryable object struct
// Always use NewItem constructor instead of &Item{}
type Item struct {
	client    *gosip.SPClient
	config    *RequestConfig
	endpoint  string
	modifiers *ODataMods
}

// GenericItemInfo - list's generic item response payload structure
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

// ItemResp - item response type with helper processor methods
type ItemResp []byte

// ListItemAllFieldsResp - item fields response type with helper processor methods
type ListItemAllFieldsResp []byte

// NewItem - Item struct constructor function
func NewItem(client *gosip.SPClient, endpoint string, config *RequestConfig) *Item {
	return &Item{
		client:    client,
		endpoint:  endpoint,
		config:    config,
		modifiers: NewODataMods(),
	}
}

// ToURL gets endpoint with modificators raw URL
func (item *Item) ToURL() string {
	return toURL(item.endpoint, item.modifiers)
}

// Get gets this Item info
func (item *Item) Get() (ItemResp, error) {
	client := NewHTTPClient(item.client)
	return client.Get(item.ToURL(), item.config)
}

// Delete deletes this Item (can't be restored from a recycle bin)
func (item *Item) Delete() error {
	client := NewHTTPClient(item.client)
	_, err := client.Delete(item.endpoint, item.config)
	return err
}

// Recycle moves this item to the recycle bin
func (item *Item) Recycle() error {
	endpoint := fmt.Sprintf("%s/Recycle", item.endpoint)
	client := NewHTTPClient(item.client)
	_, err := client.Post(endpoint, nil, item.config)
	return err
}

// Update updates item's metadata. `body` parameter is byte array representation of JSON string payload relevant to item metadata object.
func (item *Item) Update(body []byte) (ItemResp, error) {
	body = patchMetadataTypeCB(body, func() string {
		endpoint := getPriorEndpoint(item.endpoint, "/Items")
		cacheKey := strings.ToLower(endpoint + "@entitytype")
		if oDataType, found := storage.Get(cacheKey); found {
			return oDataType.(string)
		}
		list := NewList(item.client, endpoint, nil)
		oDataType, _ := list.GetEntityType()
		storage.Set(cacheKey, oDataType, 0)
		return oDataType
	})
	client := NewHTTPClient(item.client)
	return client.Update(item.endpoint, bytes.NewBuffer(body), item.config)
}

// Roles gets Roles API instance queryable collection for this Item
func (item *Item) Roles() *Roles {
	return NewRoles(item.client, item.endpoint, item.config)
}

// Attachments gets attachments collection for this Item
func (item *Item) Attachments() *Attachments {
	return NewAttachments(
		item.client,
		fmt.Sprintf("%s/AttachmentFiles", item.endpoint),
		item.config,
	)
}

// ParentList gets this Item's Lists API object
func (item *Item) ParentList() *List {
	return NewList(
		item.client,
		fmt.Sprintf("%s/ParentList", item.endpoint),
		item.config,
	)
}

// Records gets Records API instance object for this Item (inplace records manipulation)
func (item *Item) Records() *Records {
	return NewRecords(item)
}

// ContextInfo gets current context information
func (item *Item) ContextInfo() (*ContextInfo, error) {
	return NewContext(item.client, item.ToURL(), item.config).Get()
}

/* Response helpers */

// Data : to get typed data
func (itemResp *ItemResp) Data() *GenericItemInfo {
	data := NormalizeODataItem(*itemResp)
	data = fixDatesInResponse(data, []string{"Created", "Modified"})
	res := &GenericItemInfo{}
	_ = json.Unmarshal(data, &res)
	return res
}

// Normalized returns normalized body
// func (itemResp *ItemResp) Normalized() []byte {
// 	return NormalizeODataItem(*itemResp)
// }
