package api

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/koltyakov/gosip"
)

// Folders represent SharePoint Lists & Document Libraries Folders API queryable collection struct
type Folders struct {
	client    *gosip.SPClient
	config    *RequestConfig
	endpoint  string
	modifiers map[string]string
}

// FoldersResp - folders response type with helper processor methods
type FoldersResp []byte

// NewFolders - Folders struct constructor function
func NewFolders(client *gosip.SPClient, endpoint string, config *RequestConfig) *Folders {
	return &Folders{
		client:   client,
		endpoint: endpoint,
		config:   config,
	}
}

// ToURL ...
func (folders *Folders) ToURL() string {
	return folders.endpoint
}

// Conf ...
func (folders *Folders) Conf(config *RequestConfig) *Folders {
	folders.config = config
	return folders
}

// Select ...
func (folders *Folders) Select(oDataSelect string) *Folders {
	if folders.modifiers == nil {
		folders.modifiers = make(map[string]string)
	}
	folders.modifiers["$select"] = oDataSelect
	return folders
}

// Expand ...
func (folders *Folders) Expand(oDataExpand string) *Folders {
	if folders.modifiers == nil {
		folders.modifiers = make(map[string]string)
	}
	folders.modifiers["$expand"] = oDataExpand
	return folders
}

// Filter ...
func (folders *Folders) Filter(oDataFilter string) *Folders {
	if folders.modifiers == nil {
		folders.modifiers = make(map[string]string)
	}
	folders.modifiers["$filter"] = oDataFilter
	return folders
}

// Top ...
func (folders *Folders) Top(oDataTop int) *Folders {
	if folders.modifiers == nil {
		folders.modifiers = make(map[string]string)
	}
	folders.modifiers["$top"] = fmt.Sprintf("%d", oDataTop)
	return folders
}

// OrderBy ...
func (folders *Folders) OrderBy(oDataOrderBy string, ascending bool) *Folders {
	direction := "asc"
	if !ascending {
		direction = "desc"
	}
	if folders.modifiers == nil {
		folders.modifiers = make(map[string]string)
	}
	if folders.modifiers["$orderby"] != "" {
		folders.modifiers["$orderby"] += ","
	}
	folders.modifiers["$orderby"] += fmt.Sprintf("%s %s", oDataOrderBy, direction)
	return folders
}

// Get ...
func (folders *Folders) Get() (FoldersResp, error) {
	apiURL, _ := url.Parse(folders.endpoint)
	query := url.Values{}
	for k, v := range folders.modifiers {
		query.Add(k, trimMultiline(v))
	}
	apiURL.RawQuery = query.Encode()
	sp := NewHTTPClient(folders.client)
	return sp.Get(apiURL.String(), getConfHeaders(folders.config))
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
