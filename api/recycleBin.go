package api

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/koltyakov/gosip"
)

// RecycleBin represents SharePoint Recycle Bin API queryable collection struct
// Always use NewRecycleBin constructor instead of &RecycleBin{}
type RecycleBin struct {
	client    *gosip.SPClient
	config    *RequestConfig
	endpoint  string
	modifiers *ODataMods
}

// RecycledItem ...
type RecycledItem struct {
	AuthorEmail               string      `json:"AuthorEmail"`
	AuthorName                string      `json:"AuthorName"`
	DeletedByEmail            string      `json:"DeletedByEmail"`
	DeletedByName             string      `json:"DeletedByName"`
	DeletedDate               time.Time   `json:"DeletedDate"`
	DeletedDateLocalFormatted string      `json:"DeletedDateLocalFormatted"`
	DirName                   string      `json:"DirName"`
	ID                        string      `json:"Id"`
	ItemState                 int         `json:"ItemState"`
	ItemType                  int         `json:"ItemType"`
	LeafName                  string      `json:"LeafName"`
	Size                      int         `json:"Size"`
	Title                     string      `json:"Title"`
	LeafNamePath              *DecodedURL `json:"LeafNamePath"`
	DirNamePath               *DecodedURL `json:"DirNamePath"`
}

// RecycleBinResp - recycle bin response type with helper processor methods
type RecycleBinResp []byte

// NewRecycleBin - RecycleBin struct constructor function
func NewRecycleBin(client *gosip.SPClient, endpoint string, config *RequestConfig) *RecycleBin {
	return &RecycleBin{
		client:    client,
		endpoint:  endpoint,
		config:    config,
		modifiers: NewODataMods(),
	}
}

// ToURL gets endpoint with modificators raw URL
func (recycleBin *RecycleBin) ToURL() string {
	return toURL(recycleBin.endpoint, recycleBin.modifiers)
}

// Conf receives custom request config definition, e.g. custom headers, custom OData mod
func (recycleBin *RecycleBin) Conf(config *RequestConfig) *RecycleBin {
	recycleBin.config = config
	return recycleBin
}

// Select adds $select OData modifier
func (recycleBin *RecycleBin) Select(oDataSelect string) *RecycleBin {
	recycleBin.modifiers.AddSelect(oDataSelect)
	return recycleBin
}

// Expand adds $expand OData modifier
func (recycleBin *RecycleBin) Expand(oDataExpand string) *RecycleBin {
	recycleBin.modifiers.AddExpand(oDataExpand)
	return recycleBin
}

// Filter adds $filter OData modifier
func (recycleBin *RecycleBin) Filter(oDataFilter string) *RecycleBin {
	recycleBin.modifiers.AddFilter(oDataFilter)
	return recycleBin
}

// Top adds $top OData modifier
func (recycleBin *RecycleBin) Top(oDataTop int) *RecycleBin {
	recycleBin.modifiers.AddTop(oDataTop)
	return recycleBin
}

// OrderBy adds $orderby OData modifier
func (recycleBin *RecycleBin) OrderBy(oDataOrderBy string, ascending bool) *RecycleBin {
	recycleBin.modifiers.AddOrderBy(oDataOrderBy, ascending)
	return recycleBin
}

// Get gets recycled items queryable collection
func (recycleBin *RecycleBin) Get() (RecycleBinResp, error) {
	sp := NewHTTPClient(recycleBin.client)
	return sp.Get(recycleBin.ToURL(), getConfHeaders(recycleBin.config))
}

// GetByID gets a recycled item by its ID
func (recycleBin *RecycleBin) GetByID(itemID string) *RecycleBinItem {
	return NewRecycleBinItem(
		recycleBin.client,
		fmt.Sprintf("%s('%s')", recycleBin.endpoint, itemID),
		recycleBin.config,
	)
}

/* Response helpers */

// Data : to get typed data
func (recycleBinResp *RecycleBinResp) Data() []RecycleBinItemResp {
	collection, _ := normalizeODataCollection(*recycleBinResp)
	items := []RecycleBinItemResp{}
	for _, item := range collection {
		items = append(items, RecycleBinItemResp(item))
	}
	return items
}

// Normalized returns normalized body
func (recycleBinResp *RecycleBinResp) Normalized() []byte {
	normalized, _ := NormalizeODataCollection(*recycleBinResp)
	return normalized
}

/* Recycle bin item */

// RecycleBinItem represent SharePoint Recycle Bin Item API queryable object struct
// Always use NewRecycleBinItem constructor instead of &RecycleBinItem{}
type RecycleBinItem struct {
	client    *gosip.SPClient
	config    *RequestConfig
	endpoint  string
	modifiers *ODataMods
}

// RecycleBinItemResp - recycle bin item response type with helper processor methods
type RecycleBinItemResp []byte

// NewRecycleBinItem - RecycleBinItem struct constructor function
func NewRecycleBinItem(client *gosip.SPClient, endpoint string, config *RequestConfig) *RecycleBinItem {
	return &RecycleBinItem{
		client:    client,
		endpoint:  endpoint,
		config:    config,
		modifiers: NewODataMods(),
	}
}

// Get gets this recycle item data object
func (recycleBinItem *RecycleBinItem) Get() (RecycleBinItemResp, error) {
	sp := NewHTTPClient(recycleBinItem.client)
	return sp.Get(recycleBinItem.endpoint, getConfHeaders(recycleBinItem.config))
}

// Restore restores this recycled item
func (recycleBinItem *RecycleBinItem) Restore() error {
	endpoint := fmt.Sprintf("%s/Restore()", recycleBinItem.endpoint)
	sp := NewHTTPClient(recycleBinItem.client)
	_, err := sp.Post(endpoint, nil, getConfHeaders(recycleBinItem.config))
	return err
}

/* Response helpers */

// Data : to get typed data
func (recycleBinItemResp *RecycleBinItemResp) Data() *RecycledItem {
	data := NormalizeODataItem(*recycleBinItemResp)
	res := &RecycledItem{}
	json.Unmarshal(data, &res)
	return res
}

// Normalized returns normalized body
func (recycleBinItemResp *RecycleBinItemResp) Normalized() []byte {
	return NormalizeODataItem(*recycleBinItemResp)
}
