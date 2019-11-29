package api

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/koltyakov/gosip"
)

// File ...
type File struct {
	client    *gosip.SPClient
	config    *RequestConfig
	endpoint  string
	modifiers map[string]string
}

// NewFile ...
func NewFile(client *gosip.SPClient, endpoint string, config *RequestConfig) *File {
	return &File{
		client:   client,
		endpoint: endpoint,
		config:   config,
	}
}

// ToURL ...
func (file *File) ToURL() string {
	apiURL, _ := url.Parse(file.endpoint)
	query := url.Values{}
	for k, v := range file.modifiers {
		query.Add(k, trimMultiline(v))
	}
	apiURL.RawQuery = query.Encode()
	return apiURL.String()
}

// Conf ...
func (file *File) Conf(config *RequestConfig) *File {
	file.config = config
	return file
}

// Select ...
func (file *File) Select(oDataSelect string) *File {
	if file.modifiers == nil {
		file.modifiers = make(map[string]string)
	}
	file.modifiers["$select"] = oDataSelect
	return file
}

// Expand ...
func (file *File) Expand(oDataExpand string) *File {
	if file.modifiers == nil {
		file.modifiers = make(map[string]string)
	}
	file.modifiers["$expand"] = oDataExpand
	return file
}

// Get ...
func (file *File) Get() ([]byte, error) {
	sp := NewHTTPClient(file.client)
	return sp.Get(file.ToURL(), getConfHeaders(file.config))
}

// Delete ...
func (file *File) Delete() ([]byte, error) {
	sp := NewHTTPClient(file.client)
	return sp.Delete(file.endpoint, getConfHeaders(file.config))
}

// Recycle ...
func (file *File) Recycle() ([]byte, error) {
	sp := NewHTTPClient(file.client)
	endpoint := fmt.Sprintf("%s/Recycle", file.endpoint)
	return sp.Post(endpoint, nil, getConfHeaders(file.config))
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

	item := NewItem(
		file.client,
		fmt.Sprintf(
			"%s/_api/%s",
			file.client.AuthCnfg.GetSiteURL(),
			res.D.Metadata.URI,
		),
		file.config,
	)
	return item, nil
}
