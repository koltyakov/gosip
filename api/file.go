package api

import (
	"fmt"
	"net/url"

	"github.com/koltyakov/gosip"
)

// File ...
type File struct {
	client   *gosip.SPClient
	conf     *Conf
	endpoint string
	oSelect  string
	oExpand  string
}

// Conf ...
func (file *File) Conf(conf *Conf) *File {
	file.conf = conf
	return file
}

// Select ...
func (file *File) Select(oDataSelect string) *File {
	file.oSelect = oDataSelect
	return file
}

// Expand ...
func (file *File) Expand(oDataExpand string) *File {
	file.oExpand = oDataExpand
	return file
}

// Get ...
func (file *File) Get() ([]byte, error) {
	apiURL, _ := url.Parse(file.endpoint)
	query := url.Values{}
	if file.oSelect != "" {
		query.Add("$select", TrimMultiline(file.oSelect))
	}
	if file.oExpand != "" {
		query.Add("$expand", TrimMultiline(file.oExpand))
	}
	apiURL.RawQuery = query.Encode()
	sp := &HTTPClient{SPClient: file.client}
	return sp.Get(apiURL.String(), GetConfHeaders(file.conf))
}

// Delete ...
func (file *File) Delete() ([]byte, error) {
	sp := &HTTPClient{SPClient: file.client}
	return sp.Delete(file.endpoint, GetConfHeaders(file.conf))
}

// Recycle ...
func (file *File) Recycle() ([]byte, error) {
	sp := &HTTPClient{SPClient: file.client}
	endpoint := fmt.Sprintf("%s/Recycle", file.endpoint)
	return sp.Post(endpoint, nil, GetConfHeaders(file.conf))
}
