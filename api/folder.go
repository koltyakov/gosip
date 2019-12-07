package api

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/koltyakov/gosip"
)

// Folder ...
type Folder struct {
	client    *gosip.SPClient
	config    *RequestConfig
	endpoint  string
	modifiers map[string]string
}

// NewFolder ...
func NewFolder(client *gosip.SPClient, endpoint string, config *RequestConfig) *Folder {
	return &Folder{
		client:   client,
		endpoint: endpoint,
		config:   config,
	}
}

// ToURL ...
func (folder *Folder) ToURL() string {
	apiURL, _ := url.Parse(folder.endpoint)
	query := apiURL.Query() // url.Values{}
	for k, v := range folder.modifiers {
		query.Set(k, trimMultiline(v))
	}
	apiURL.RawQuery = query.Encode()
	return apiURL.String()
}

// Conf ...
func (folder *Folder) Conf(config *RequestConfig) *Folder {
	folder.config = config
	return folder
}

// Select ...
func (folder *Folder) Select(oDataSelect string) *Folder {
	if folder.modifiers == nil {
		folder.modifiers = make(map[string]string)
	}
	folder.modifiers["$select"] = oDataSelect
	return folder
}

// Expand ...
func (folder *Folder) Expand(oDataExpand string) *Folder {
	if folder.modifiers == nil {
		folder.modifiers = make(map[string]string)
	}
	folder.modifiers["$expand"] = oDataExpand
	return folder
}

// Get ...
func (folder *Folder) Get() ([]byte, error) {
	sp := NewHTTPClient(folder.client)
	return sp.Get(folder.ToURL(), getConfHeaders(folder.config))
}

// Delete ...
func (folder *Folder) Delete() ([]byte, error) {
	sp := NewHTTPClient(folder.client)
	return sp.Delete(folder.endpoint, getConfHeaders(folder.config))
}

// Recycle ...
func (folder *Folder) Recycle() ([]byte, error) {
	sp := NewHTTPClient(folder.client)
	endpoint := fmt.Sprintf("%s/Recycle", folder.endpoint)
	return sp.Post(endpoint, nil, getConfHeaders(folder.config))
}

// Folders ...
func (folder *Folder) Folders() *Folders {
	return NewFolders(
		folder.client,
		fmt.Sprintf("%s/Folders", folder.endpoint),
		folder.config,
	)
}

// Files ...
func (folder *Folder) Files() *Files {
	return NewFiles(
		folder.client,
		fmt.Sprintf("%s/Files", folder.endpoint),
		folder.config,
	)
}

// GetItem ...
func (folder *Folder) GetItem() (*Item, error) {
	endpoint := fmt.Sprintf("%s/ListItemAllFields", folder.endpoint)
	apiURL, _ := url.Parse(endpoint)
	query := url.Values{}
	query.Add("$select", "Id")
	apiURL.RawQuery = query.Encode()
	sp := NewHTTPClient(folder.client)

	data, err := sp.Get(apiURL.String(), HeadersPresets.Verbose.Headers)
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
		folder.client,
		fmt.Sprintf(
			"%s/_api/%s",
			folder.client.AuthCnfg.GetSiteURL(),
			res.D.Metadata.URI,
		),
		folder.config,
	)
	return item, nil
}
