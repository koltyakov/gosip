package api

import (
	"fmt"
	"net/url"

	"github.com/koltyakov/gosip"
)

// Folders ...
type Folders struct {
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
func (folders *Folders) Conf(config *RequestConfig) *Folders {
	folders.config = config
	return folders
}

// Select ...
func (folders *Folders) Select(oDataSelect string) *Folders {
	folders.oSelect = oDataSelect
	return folders
}

// Expand ...
func (folders *Folders) Expand(oDataExpand string) *Folders {
	folders.oExpand = oDataExpand
	return folders
}

// Filter ...
func (folders *Folders) Filter(oDataFilter string) *Folders {
	folders.oFilter = oDataFilter
	return folders
}

// Top ...
func (folders *Folders) Top(oDataTop int) *Folders {
	folders.oTop = oDataTop
	return folders
}

// OrderBy ...
func (folders *Folders) OrderBy(oDataOrderBy string, ascending bool) *Folders {
	direction := "asc"
	if !ascending {
		direction = "desc"
	}
	if folders.oOrderBy != "" {
		folders.oOrderBy += ","
	}
	folders.oOrderBy += fmt.Sprintf("%s %s", oDataOrderBy, direction)
	return folders
}

// Get ...
func (folders *Folders) Get() ([]byte, error) {
	apiURL, _ := url.Parse(folders.endpoint)
	query := url.Values{}
	if folders.oSelect != "" {
		query.Add("$select", trimMultiline(folders.oSelect))
	}
	if folders.oExpand != "" {
		query.Add("$expand", trimMultiline(folders.oExpand))
	}
	if folders.oFilter != "" {
		query.Add("$filter", trimMultiline(folders.oFilter))
	}
	if folders.oTop != 0 {
		query.Add("$top", fmt.Sprintf("%d", folders.oTop))
	}
	if folders.oOrderBy != "" {
		query.Add("$orderBy", trimMultiline(folders.oOrderBy))
	}
	apiURL.RawQuery = query.Encode()
	sp := NewHTTPClient(folders.client)
	return sp.Get(apiURL.String(), getConfHeaders(folders.config))
}

// Add ...
func (folders *Folders) Add(folderName string) ([]byte, error) {
	sp := NewHTTPClient(folders.client)
	endpoint := fmt.Sprintf("%s/Add('%s')", folders.endpoint, folderName)
	return sp.Post(endpoint, nil, getConfHeaders(folders.config))
}

// GetByName ...
func (folders *Folders) GetByName(folderName string) *Folder {
	return &Folder{
		client: folders.client,
		config: folders.config,
		endpoint: fmt.Sprintf("%s('%s')",
			folders.endpoint,
			folderName,
		),
	}
}
