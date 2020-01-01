package api

import (
	"encoding/json"
	"fmt"

	"github.com/koltyakov/gosip"
)

// Folders represent SharePoint Lists & Document Libraries Folders API queryable collection struct
// Always use NewFolders constructor instead of &Folders{}
type Folders struct {
	client    *gosip.SPClient
	config    *RequestConfig
	endpoint  string
	modifiers *ODataMods
}

// FoldersResp - folders response type with helper processor methods
type FoldersResp []byte

// NewFolders - Folders struct constructor function
func NewFolders(client *gosip.SPClient, endpoint string, config *RequestConfig) *Folders {
	return &Folders{
		client:    client,
		endpoint:  endpoint,
		config:    config,
		modifiers: NewODataMods(),
	}
}

// ToURL gets endpoint with modificators raw URL
func (folders *Folders) ToURL() string {
	// return folders.endpoint
	return toURL(folders.endpoint, folders.modifiers)
}

// Conf receives custom request config definition, e.g. custom headers, custom OData mod
func (folders *Folders) Conf(config *RequestConfig) *Folders {
	folders.config = config
	return folders
}

// Select ...
func (folders *Folders) Select(oDataSelect string) *Folders {
	folders.modifiers.AddSelect(oDataSelect)
	return folders
}

// Expand ...
func (folders *Folders) Expand(oDataExpand string) *Folders {
	folders.modifiers.AddExpand(oDataExpand)
	return folders
}

// Filter ...
func (folders *Folders) Filter(oDataFilter string) *Folders {
	folders.modifiers.AddFilter(oDataFilter)
	return folders
}

// Top ...
func (folders *Folders) Top(oDataTop int) *Folders {
	folders.modifiers.AddTop(oDataTop)
	return folders
}

// OrderBy ...
func (folders *Folders) OrderBy(oDataOrderBy string, ascending bool) *Folders {
	folders.modifiers.AddOrderBy(oDataOrderBy, ascending)
	return folders
}

// Get ...
func (folders *Folders) Get() (FoldersResp, error) {
	sp := NewHTTPClient(folders.client)
	return sp.Get(folders.ToURL(), getConfHeaders(folders.config))
}

// Add ...
func (folders *Folders) Add(folderName string) (FolderResp, error) {
	sp := NewHTTPClient(folders.client)
	endpoint := fmt.Sprintf("%s/Add('%s')", folders.endpoint, folderName)
	return sp.Post(endpoint, nil, getConfHeaders(folders.config))
}

// GetByName ...
func (folders *Folders) GetByName(folderName string) *Folder {
	return NewFolder(
		folders.client,
		fmt.Sprintf("%s('%s')", folders.endpoint, folderName),
		folders.config,
	)
}

/* Response helpers */

// Data : to get typed data
func (foldersResp *FoldersResp) Data() []FolderResp {
	collection, _ := parseODataCollection(*foldersResp)
	folders := []FolderResp{}
	for _, ct := range collection {
		folders = append(folders, FolderResp(ct))
	}
	return folders
}

// Unmarshal : to unmarshal to custom object
func (foldersResp *FoldersResp) Unmarshal(obj interface{}) error {
	data, _ := parseODataCollectionPlain(*foldersResp)
	return json.Unmarshal(data, obj)
}
