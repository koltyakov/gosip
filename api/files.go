package api

import (
	"fmt"
	"net/url"

	"github.com/koltyakov/gosip"
)

// Files ...
type Files struct {
	client    *gosip.SPClient
	config    *RequestConfig
	endpoint  string
	modifiers map[string]string
}

// NewFiles ...
func NewFiles(client *gosip.SPClient, endpoint string, config *RequestConfig) *Files {
	return &Files{
		client:   client,
		endpoint: endpoint,
		config:   config,
	}
}

// ToURL ...
func (files *Files) ToURL() string {
	apiURL, _ := url.Parse(files.endpoint)
	query := url.Values{}
	for k, v := range files.modifiers {
		query.Add(k, trimMultiline(v))
	}
	apiURL.RawQuery = query.Encode()
	return apiURL.String()
}

// Conf ...
func (files *Files) Conf(config *RequestConfig) *Files {
	files.config = config
	return files
}

// Select ...
func (files *Files) Select(oDataSelect string) *Files {
	if files.modifiers == nil {
		files.modifiers = make(map[string]string)
	}
	files.modifiers["$select"] = oDataSelect
	return files
}

// Expand ...
func (files *Files) Expand(oDataExpand string) *Files {
	if files.modifiers == nil {
		files.modifiers = make(map[string]string)
	}
	files.modifiers["$expand"] = oDataExpand
	return files
}

// Filter ...
func (files *Files) Filter(oDataFilter string) *Files {
	if files.modifiers == nil {
		files.modifiers = make(map[string]string)
	}
	files.modifiers["$filter"] = oDataFilter
	return files
}

// Top ...
func (files *Files) Top(oDataTop int) *Files {
	if files.modifiers == nil {
		files.modifiers = make(map[string]string)
	}
	files.modifiers["$top"] = fmt.Sprintf("%d", oDataTop)
	return files
}

// OrderBy ...
func (files *Files) OrderBy(oDataOrderBy string, ascending bool) *Files {
	direction := "asc"
	if !ascending {
		direction = "desc"
	}
	if files.modifiers == nil {
		files.modifiers = make(map[string]string)
	}
	if files.modifiers["$orderby"] != "" {
		files.modifiers["$orderby"] += ","
	}
	files.modifiers["$orderby"] += fmt.Sprintf("%s %s", oDataOrderBy, direction)
	return files
}

// Get ...
func (files *Files) Get() ([]byte, error) {
	sp := NewHTTPClient(files.client)
	return sp.Get(files.ToURL(), getConfHeaders(files.config))
}

// GetByName ...
func (files *Files) GetByName(fileName string) *File {
	return NewFile(
		files.client,
		fmt.Sprintf("%s('%s')", files.endpoint, fileName),
		files.config,
	)
}

// Add ...
func (files *Files) Add(name string, content []byte, overwrite bool) ([]byte, error) {
	sp := NewHTTPClient(files.client)
	endpoint := fmt.Sprintf("%s/Add(overwrite=%t,url='%s')", files.endpoint, overwrite, name)
	return sp.Post(endpoint, content, getConfHeaders(files.config))
}
