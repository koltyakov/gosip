package api

import (
	"encoding/json"
	"fmt"
	"net/url"
	"time"

	"github.com/koltyakov/gosip"
)

// Folder represents SharePoint Lists & Document Libraries Folder API queryable object struct
type Folder struct {
	client    *gosip.SPClient
	config    *RequestConfig
	endpoint  string
	modifiers map[string]string
}

// FolderInfo - folder API response payload structure
type FolderInfo struct {
	Exists            bool      `json:"Exists"`
	IsWOPIEnabled     bool      `json:"IsWOPIEnabled"`
	ItemCount         int       `json:"ItemCount"`
	Name              string    `json:"Name"`
	ProgID            string    `json:"ProgID"`
	ServerRelativeURL string    `json:"ServerRelativeUrl"`
	TimeCreated       time.Time `json:"TimeCreated"`
	TimeLastModified  time.Time `json:"TimeLastModified"`
	UniqueID          string    `json:"UniqueId"`
	WelcomePage       string    `json:"WelcomePage"`
}

// FolderResp - folder response type with helper processor methods
type FolderResp []byte

// NewFolder ...
func NewFolder(client *gosip.SPClient, endpoint string, config *RequestConfig) *Folder {
	return &Folder{
		client:   client,
		endpoint: endpoint,
		config:   config,
	}
}

// ToURL gets endpoint with modificators raw URL ...
func (folder *Folder) ToURL() string {
	return toURL(folder.endpoint, folder.modifiers)
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
func (folder *Folder) Get() (FolderResp, error) {
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

// ContextInfo ...
func (folder *Folder) ContextInfo() (*ContextInfo, error) {
	return NewContext(folder.client, folder.ToURL(), folder.config).Get()
}

/* Response helpers */

// Data : to get typed data
func (folderResp *FolderResp) Data() *FolderInfo {
	data := parseODataItem(*folderResp)
	res := &FolderInfo{}
	json.Unmarshal(data, &res)
	return res
}

// Unmarshal : to unmarshal to custom object
func (folderResp *FolderResp) Unmarshal(obj interface{}) error {
	data := parseODataItem(*folderResp)
	return json.Unmarshal(data, obj)
}
