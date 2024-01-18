package api

import (
	"context"
	"fmt"
	"time"

	"github.com/koltyakov/gosip"
)

//go:generate ggen -ent RecycleBin -item RecycleBinItem -conf -coll -mods Select,Expand,Filter,Top,OrderBy -helpers Data,Normalized
//go:generate ggen -ent RecycleBinItem -helpers Data,Normalized

// RecycleBin represents SharePoint Recycle Bin API queryable collection struct
// Always use NewRecycleBin constructor instead of &RecycleBin{}
type RecycleBin struct {
	client    *gosip.SPClient
	config    *RequestConfig
	endpoint  string
	modifiers *ODataMods
}

// RecycleBinItemInfo ...
type RecycleBinItemInfo struct {
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

// Get gets recycled items queryable collection
func (recycleBin *RecycleBin) Get(ctx context.Context) (RecycleBinResp, error) {
	client := NewHTTPClient(recycleBin.client)
	return client.Get(ctx, recycleBin.ToURL(), recycleBin.config)
}

// GetByID gets a recycled item by its ID
func (recycleBin *RecycleBin) GetByID(itemID string) *RecycleBinItem {
	return NewRecycleBinItem(
		recycleBin.client,
		fmt.Sprintf("%s('%s')", recycleBin.endpoint, itemID),
		recycleBin.config,
	)
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
func (recycleBinItem *RecycleBinItem) Get(ctx context.Context) (RecycleBinItemResp, error) {
	client := NewHTTPClient(recycleBinItem.client)
	return client.Get(ctx, recycleBinItem.endpoint, recycleBinItem.config)
}

// Restore restores this recycled item
func (recycleBinItem *RecycleBinItem) Restore(ctx context.Context) error {
	endpoint := fmt.Sprintf("%s/Restore()", recycleBinItem.endpoint)
	client := NewHTTPClient(recycleBinItem.client)
	_, err := client.Post(ctx, endpoint, nil, recycleBinItem.config)
	return err
}
