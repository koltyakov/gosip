package api

import (
	"fmt"
	"net/url"

	"github.com/koltyakov/gosip"
)

// Files ...
type Files struct {
	client   *gosip.SPClient
	conf     *Conf
	endpoint string
	oSelect  string
	oExpand  string
	oFilter  string
	oTop     int
	oOrderBy string
}

// Conf ...
func (files *Files) Conf(conf *Conf) *Files {
	files.conf = conf
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
		query.Add("$select", TrimMultiline(files.oSelect))
	}
	if files.oExpand != "" {
		query.Add("$expand", TrimMultiline(files.oExpand))
	}
	if files.oFilter != "" {
		query.Add("$filter", TrimMultiline(files.oFilter))
	}
	if files.oTop != 0 {
		query.Add("$top", fmt.Sprintf("%d", files.oTop))
	}
	if files.oOrderBy != "" {
		query.Add("$orderBy", TrimMultiline(files.oOrderBy))
	}
	apiURL.RawQuery = query.Encode()
	sp := &HTTPClient{SPClient: files.client}
	return sp.Get(apiURL.String(), GetConfHeaders(files.conf))
}

// GetByName ...
func (files *Files) GetByName(fileName string) *File {
	return &File{
		client: files.client,
		conf:   files.conf,
		endpoint: fmt.Sprintf("%s('%s')",
			files.endpoint,
			fileName,
		),
	}
}

// Add ...
func (files *Files) Add(name string, content []byte, overwrite bool) ([]byte, error) {
	sp := &HTTPClient{SPClient: files.client}
	endpoint := fmt.Sprintf("%s/Add(overwrite=%t,url='%s')", files.endpoint, overwrite, name)
	return sp.Post(endpoint, content, GetConfHeaders(files.conf))
}
