package api

import (
	"fmt"
	"net/url"

	"github.com/koltyakov/gosip"
)

// Files ...
type Files struct {
	client   *gosip.SPClient
	config   *RequestConfig
	endpoint string
	oSelect  string
	oExpand  string
	oFilter  string
	oTop     int
	oOrderBy string
}

// Conf ...
func (files *Files) Conf(config *RequestConfig) *Files {
	files.config = config
	return files
}

// Select ...
func (files *Files) Select(oDataSelect string) *Files {
	files.oSelect = oDataSelect
	return files
}

// Expand ...
func (files *Files) Expand(oDataExpand string) *Files {
	files.oExpand = oDataExpand
	return files
}

// Filter ...
func (files *Files) Filter(oDataFilter string) *Files {
	files.oFilter = oDataFilter
	return files
}

// Top ...
func (files *Files) Top(oDataTop int) *Files {
	files.oTop = oDataTop
	return files
}

// OrderBy ...
func (files *Files) OrderBy(oDataOrderBy string, ascending bool) *Files {
	direction := "asc"
	if !ascending {
		direction = "desc"
	}
	if files.oOrderBy != "" {
		files.oOrderBy += ","
	}
	files.oOrderBy += fmt.Sprintf("%s %s", oDataOrderBy, direction)
	return files
}

// Get ...
func (files *Files) Get() ([]byte, error) {
	apiURL, _ := url.Parse(files.endpoint)
	query := url.Values{}
	if files.oSelect != "" {
		query.Add("$select", trimMultiline(files.oSelect))
	}
	if files.oExpand != "" {
		query.Add("$expand", trimMultiline(files.oExpand))
	}
	if files.oFilter != "" {
		query.Add("$filter", trimMultiline(files.oFilter))
	}
	if files.oTop != 0 {
		query.Add("$top", fmt.Sprintf("%d", files.oTop))
	}
	if files.oOrderBy != "" {
		query.Add("$orderBy", trimMultiline(files.oOrderBy))
	}
	apiURL.RawQuery = query.Encode()
	sp := NewHTTPClient(files.client)
	return sp.Get(apiURL.String(), getConfHeaders(files.config))
}

// GetByName ...
func (files *Files) GetByName(fileName string) *File {
	return &File{
		client: files.client,
		config: files.config,
		endpoint: fmt.Sprintf("%s('%s')",
			files.endpoint,
			fileName,
		),
	}
}

// Add ...
func (files *Files) Add(name string, content []byte, overwrite bool) ([]byte, error) {
	sp := NewHTTPClient(files.client)
	endpoint := fmt.Sprintf("%s/Add(overwrite=%t,url='%s')", files.endpoint, overwrite, name)
	return sp.Post(endpoint, content, getConfHeaders(files.config))
}
