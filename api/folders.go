package api

import (
	"fmt"

	"github.com/koltyakov/gosip"
)

//go:generate ggen -ent Folders -conf -mods Select,Expand,Filter,Top,OrderBy

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

// Get gets folders collection response in this folder
func (folders *Folders) Get() (FoldersResp, error) {
	sp := NewHTTPClient(folders.client)
	return sp.Get(folders.ToURL(), getConfHeaders(folders.config))
}

// Add created a folder with specified name in this folder
func (folders *Folders) Add(folderName string) (FolderResp, error) {
	sp := NewHTTPClient(folders.client)
	endpoint := fmt.Sprintf("%s/Add('%s')", folders.endpoint, folderName)
	return sp.Post(endpoint, nil, getConfHeaders(folders.config))
}

// GetByName gets a folder by its name in this folder
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
	collection, _ := normalizeODataCollection(*foldersResp)
	folders := []FolderResp{}
	for _, ct := range collection {
		folders = append(folders, FolderResp(ct))
	}
	return folders
}

// Normalized returns normalized body
func (foldersResp *FoldersResp) Normalized() []byte {
	normalized, _ := NormalizeODataCollection(*foldersResp)
	return normalized
}
