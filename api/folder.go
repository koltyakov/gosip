package api

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/koltyakov/gosip"
)

// Folder ...
type Folder struct {
	client   *gosip.SPClient
	config   *RequestConfig
	endpoint string
	oSelect  string
	oExpand  string
}

// Conf ...
func (folder *Folder) Conf(config *RequestConfig) *Folder {
	folder.config = config
	return folder
}

// Select ...
func (folder *Folder) Select(oDataSelect string) *Folder {
	folder.oSelect = oDataSelect
	return folder
}

// Expand ...
func (folder *Folder) Expand(oDataExpand string) *Folder {
	folder.oExpand = oDataExpand
	return folder
}

// Get ...
func (folder *Folder) Get() ([]byte, error) {
	apiURL, _ := url.Parse(folder.endpoint)
	query := url.Values{}
	if folder.oSelect != "" {
		query.Add("$select", trimMultiline(folder.oSelect))
	}
	if folder.oExpand != "" {
		query.Add("$expand", trimMultiline(folder.oExpand))
	}
	apiURL.RawQuery = query.Encode()
	sp := NewHTTPClient(folder.client)
	return sp.Get(apiURL.String(), getConfHeaders(folder.config))
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
	return &Folders{
		client: folder.client,
		config: folder.config,
		endpoint: fmt.Sprintf(
			"%s/Folders",
			folder.endpoint,
		),
	}
}

// Files ...
func (folder *Folder) Files() *Files {
	return &Files{
		client: folder.client,
		config: folder.config,
		endpoint: fmt.Sprintf(
			"%s/Files",
			folder.endpoint,
		),
	}
}

// GetItem ...
func (folder *Folder) GetItem() (*Item, error) {
	endpoint := fmt.Sprintf("%s/ListItemAllFields", folder.endpoint)
	apiURL, _ := url.Parse(endpoint)
	query := url.Values{}
	query.Add("$select", "Id")
	apiURL.RawQuery = query.Encode()
	sp := NewHTTPClient(folder.client)

	headers := make(map[string]string)
	headers["Accept"] = "application/json;odata=verbose"
	headers["Content-Type"] = "application/json;odata=verbose;charset=utf-8"

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
		client: folder.client,
		config: folder.config,
		endpoint: fmt.Sprintf("%s/_api/%s",
			folder.client.AuthCnfg.GetSiteURL(),
			res.D.Metadata.URI,
		),
	}, nil
}
