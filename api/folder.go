package api

import (
	"encoding/json"
	"fmt"
	"net/url"
	"time"

	"github.com/koltyakov/gosip"
)

// Folder represents SharePoint Lists & Document Libraries Folder API queryable object struct
// Always use NewFolder constructor instead of &Folder{}
type Folder struct {
	client    *gosip.SPClient
	config    *RequestConfig
	endpoint  string
	modifiers *ODataMods
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
		client:    client,
		endpoint:  endpoint,
		config:    config,
		modifiers: NewODataMods(),
	}
}

// ToURL gets endpoint with modificators raw URL
func (folder *Folder) ToURL() string {
	return toURL(folder.endpoint, folder.modifiers)
}

// Conf receives custom request config definition, e.g. custom headers, custom OData mod
func (folder *Folder) Conf(config *RequestConfig) *Folder {
	folder.config = config
	return folder
}

// Select adds $select OData modifier
func (folder *Folder) Select(oDataSelect string) *Folder {
	folder.modifiers.AddSelect(oDataSelect)
	return folder
}

// Expand adds $expand OData modifier
func (folder *Folder) Expand(oDataExpand string) *Folder {
	folder.modifiers.AddExpand(oDataExpand)
	return folder
}

// Get gets this folder data object
func (folder *Folder) Get() (FolderResp, error) {
	sp := NewHTTPClient(folder.client)
	return sp.Get(folder.ToURL(), getConfHeaders(folder.config))
}

// Update updates Folder's metadata with properties provided in `body` parameter
// where `body` is byte array representation of JSON string payload relevalt to SP.Folder object
func (folder *Folder) Update(body []byte) (FolderResp, error) {
	body = patchMetadataType(body, "SP.Folder")
	sp := NewHTTPClient(folder.client)
	return sp.Update(folder.endpoint, body, getConfHeaders(folder.config))
}

// Delete deletes this folder (can't be restored from a recycle bin)
func (folder *Folder) Delete() error {
	sp := NewHTTPClient(folder.client)
	_, err := sp.Delete(folder.endpoint, getConfHeaders(folder.config))
	return err
}

// Recycle moves this folder to the recycle bin
func (folder *Folder) Recycle() error {
	sp := NewHTTPClient(folder.client)
	endpoint := fmt.Sprintf("%s/Recycle", folder.endpoint)
	_, err := sp.Post(endpoint, nil, getConfHeaders(folder.config))
	return err
}

// Folders gets subfolders queryable collection
func (folder *Folder) Folders() *Folders {
	return NewFolders(
		folder.client,
		fmt.Sprintf("%s/Folders", folder.endpoint),
		folder.config,
	)
}

// ParentFolder gets parent folder of this folder
func (folder *Folder) ParentFolder() *Folder {
	return NewFolder(
		folder.client,
		fmt.Sprintf("%s/ParentFolder", folder.endpoint),
		folder.config,
	)
}

// Props gets Properties API instance queryable collection for this Folder
func (folder *Folder) Props() *Properties {
	return NewProperties(
		folder.client,
		fmt.Sprintf("%s/Properties", folder.endpoint),
		folder.config,
	)
}

// Files gets files queryable collection in this folder
func (folder *Folder) Files() *Files {
	return NewFiles(
		folder.client,
		fmt.Sprintf("%s/Files", folder.endpoint),
		folder.config,
	)
}

// ListItemAllFields gets this folder Item data object metadata
func (folder *Folder) ListItemAllFields() (ListItemAllFieldsResp, error) {
	endpoint := fmt.Sprintf("%s/ListItemAllFields", folder.endpoint)
	apiURL, _ := url.Parse(endpoint)

	query := apiURL.Query()
	for k, v := range folder.modifiers.Get() {
		query.Set(k, TrimMultiline(v))
	}

	apiURL.RawQuery = query.Encode()
	sp := NewHTTPClient(folder.client)

	data, err := sp.Get(apiURL.String(), getConfHeaders(folder.config))
	if err != nil {
		return nil, err
	}
	data = NormalizeODataItem(data)
	return data, nil
}

// GetItem gets this folder Item API object metadata
func (folder *Folder) GetItem() (*Item, error) {
	scoped := NewFolder(folder.client, folder.endpoint, folder.config)
	data, err := scoped.Conf(HeadersPresets.Verbose).Select("Id").ListItemAllFields()
	if err != nil {
		return nil, err
	}

	res := &struct {
		Metadata struct {
			URI string `json:"uri"`
		} `json:"__metadata"`
	}{}

	err = json.Unmarshal(data, &res)
	if err != nil {
		return nil, err
	}

	item := NewItem(
		folder.client,
		res.Metadata.URI,
		folder.config,
	)
	return item, nil
}

// ContextInfo gets current Context Info object data
func (folder *Folder) ContextInfo() (*ContextInfo, error) {
	return NewContext(folder.client, folder.ToURL(), folder.config).Get()
}

// ToDo:
// StorageMetrics

/* Response helpers */

// Data : to get typed data
func (folderResp *FolderResp) Data() *FolderInfo {
	data := NormalizeODataItem(*folderResp)
	res := &FolderInfo{}
	json.Unmarshal(data, &res)
	return res
}

// Unmarshal : to unmarshal to custom object
func (folderResp *FolderResp) Unmarshal(obj interface{}) error {
	data := NormalizeODataItem(*folderResp)
	return json.Unmarshal(data, obj)
}
