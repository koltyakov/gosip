package api

import (
	"fmt"

	"github.com/koltyakov/gosip"
)

//go:generate ggen -ent Files -conf -mods Select,Expand,Filter,Top,OrderBy

// Files represent SharePoint Files API queryable collection struct
// Always use NewFiles constructor instead of &Files{}
type Files struct {
	client    *gosip.SPClient
	config    *RequestConfig
	endpoint  string
	modifiers *ODataMods
}

// FilesResp - files response type with helper processor methods
type FilesResp []byte

// NewFiles - Files struct constructor function
func NewFiles(client *gosip.SPClient, endpoint string, config *RequestConfig) *Files {
	return &Files{
		client:    client,
		endpoint:  endpoint,
		config:    config,
		modifiers: NewODataMods(),
	}
}

// ToURL gets endpoint with modificators raw URL
func (files *Files) ToURL() string {
	return toURL(files.endpoint, files.modifiers)
}

// Get gets Files collection response
func (files *Files) Get() (FilesResp, error) {
	sp := NewHTTPClient(files.client)
	return sp.Get(files.ToURL(), getConfHeaders(files.config))
}

// GetByName gets a file by its name
func (files *Files) GetByName(fileName string) *File {
	return NewFile(
		files.client,
		fmt.Sprintf("%s('%s')", files.endpoint, fileName),
		files.config,
	)
}

// Add uploads file into the folder
func (files *Files) Add(name string, content []byte, overwrite bool) (FileResp, error) {
	sp := NewHTTPClient(files.client)
	endpoint := fmt.Sprintf("%s/Add(overwrite=%t,url='%s')", files.endpoint, overwrite, name)
	return sp.Post(endpoint, content, getConfHeaders(files.config))
}

/* Response helpers */

// Data : to get typed data
func (filesResp *FilesResp) Data() []FileResp {
	collection, _ := normalizeODataCollection(*filesResp)
	files := []FileResp{}
	for _, ct := range collection {
		files = append(files, FileResp(ct))
	}
	return files
}

// Normalized returns normalized body
func (filesResp *FilesResp) Normalized() []byte {
	normalized, _ := NormalizeODataCollection(*filesResp)
	return normalized
}
