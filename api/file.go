package api

import (
	"encoding/json"
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
		query.Add("$select", trimMultiline(file.oSelect))
	}
	if file.oExpand != "" {
		query.Add("$expand", trimMultiline(file.oExpand))
	}
	apiURL.RawQuery = query.Encode()
	sp := NewHTTPClient(file.client)
	return sp.Get(apiURL.String(), getConfHeaders(file.conf))
}

// Delete ...
func (file *File) Delete() ([]byte, error) {
	sp := NewHTTPClient(file.client)
	return sp.Delete(file.endpoint, getConfHeaders(file.conf))
}

// Recycle ...
func (file *File) Recycle() ([]byte, error) {
	sp := NewHTTPClient(file.client)
	endpoint := fmt.Sprintf("%s/Recycle", file.endpoint)
	return sp.Post(endpoint, nil, getConfHeaders(file.conf))
}

// GetItem ...
func (file *File) GetItem() (*Item, error) {
	endpoint := fmt.Sprintf("%s/ListItemAllFields", file.endpoint)
	apiURL, _ := url.Parse(endpoint)
	query := url.Values{}
	query.Add("$select", "Id")
	apiURL.RawQuery = query.Encode()
	sp := NewHTTPClient(file.client)

	headers := map[string]string{
		"Accept":       "application/json;odata=verbose",
		"Content-Type": "application/json;odata=verbose;charset=utf-8",
	}

	data, err := sp.Get(apiURL.String(), headers)
	if err != nil {
		return nil, err
	}

	res := &struct {
		D struct {
			Metadata struct {
				URI string `json:"id"`
			} `json:"__metadata"`
		} `json:"d"`
	}{}

	err = json.Unmarshal(data, &res)
	if err != nil {
		return nil, err
	}

	return &Item{
		client: file.client,
		conf:   file.conf,
		endpoint: fmt.Sprintf("%s/_api/%s",
			file.client.AuthCnfg.GetSiteURL(),
			res.D.Metadata.URI,
		),
	}, nil
}
