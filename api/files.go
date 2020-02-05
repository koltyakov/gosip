package api

import (
	"bytes"
	"fmt"

	"github.com/koltyakov/gosip"
)

//go:generate ggen -ent Files -item File -conf -coll -mods Select,Expand,Filter,Top,OrderBy -helpers Data,Normalized

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
	client := NewHTTPClient(files.client)
	return client.Get(files.ToURL(), getConfHeaders(files.config))
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
	client := NewHTTPClient(files.client)
	endpoint := fmt.Sprintf("%s/Add(overwrite=%t,url='%s')", files.endpoint, overwrite, name)
	return client.Post(endpoint, bytes.NewBuffer(content), getConfHeaders(files.config))
}
