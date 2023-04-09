package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/koltyakov/gosip"
)

//go:generate ggen -ent File -conf -mods Select,Expand -helpers Data,Normalized

// File represents SharePoint File API queryable object struct
// Always use NewFile constructor instead of &File{}
type File struct {
	client    *gosip.SPClient
	config    *RequestConfig
	endpoint  string
	modifiers *ODataMods
}

// FileInfo - file API response payload structure
type FileInfo struct {
	CheckInComment       string    `json:"CheckInComment"`
	CheckOutType         int       `json:"CheckOutType"`
	ContentTag           string    `json:"ContentTag"`
	CustomizedPageStatus int       `json:"CustomizedPageStatus"`
	ETag                 string    `json:"ETag"`
	Exists               bool      `json:"Exists"`
	IrmEnabled           bool      `json:"IrmEnabled"`
	Length               int       `json:"Length,string"`
	Level                int       `json:"Level"`
	LinkingURI           string    `json:"LinkingUri"`
	LinkingURL           string    `json:"LinkingUrl"`
	MajorVersion         int       `json:"MajorVersion"`
	MinorVersion         int       `json:"MinorVersion"`
	Name                 string    `json:"Name"`
	ServerRelativeURL    string    `json:"ServerRelativeUrl"`
	TimeCreated          time.Time `json:"TimeCreated"`
	TimeLastModified     time.Time `json:"TimeLastModified"`
	Title                string    `json:"Title"`
	UIVersion            int       `json:"UIVersion"`
	UIVersionLabel       string    `json:"UIVersionLabel"`
	UniqueID             string    `json:"UniqueId"`
}

type checkInTypes struct {
	Minor     int
	Major     int
	Overwrite int
}

// CheckInTypes - available check in types
var CheckInTypes = func() *checkInTypes {
	return &checkInTypes{
		Minor:     0,
		Major:     1,
		Overwrite: 2,
	}
}()

// FileResp - file response type with helper processor methods
type FileResp []byte

// NewFile - File struct constructor function
func NewFile(client *gosip.SPClient, endpoint string, config *RequestConfig) *File {
	return &File{
		client:    client,
		endpoint:  endpoint,
		config:    config,
		modifiers: NewODataMods(),
	}
}

// ToURL gets endpoint with modificators raw URL
func (file *File) ToURL() string {
	return toURL(file.endpoint, file.modifiers)
}

// Get gets file data object
func (file *File) Get() (FileResp, error) {
	client := NewHTTPClient(file.client)
	return client.Get(file.ToURL(), file.config)
}

// Delete deletes this file skipping recycle bin
func (file *File) Delete() error {
	client := NewHTTPClient(file.client)
	_, err := client.Delete(file.endpoint, file.config)
	return err
}

// Recycle moves this file to the recycle bin
func (file *File) Recycle() error {
	client := NewHTTPClient(file.client)
	endpoint := fmt.Sprintf("%s/Recycle", file.endpoint)
	_, err := client.Post(endpoint, nil, file.config)
	return err
}

// ListItemAllFields gets this file Item data object metadata
func (file *File) ListItemAllFields() (ListItemAllFieldsResp, error) {
	endpoint := fmt.Sprintf("%s/ListItemAllFields", file.endpoint)
	apiURL, _ := url.Parse(endpoint)

	query := apiURL.Query()
	for k, v := range file.modifiers.Get() {
		query.Set(k, TrimMultiline(v))
	}

	apiURL.RawQuery = query.Encode()
	client := NewHTTPClient(file.client)

	data, err := client.Get(apiURL.String(), file.config)
	if err != nil {
		return nil, err
	}
	data = NormalizeODataItem(data)
	return data, nil
}

// GetItem gets this folder Item API object metadata
func (file *File) GetItem() (*Item, error) {
	scoped := NewFile(file.client, file.endpoint, file.config)
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
		file.client,
		res.Metadata.URI,
		file.config,
	)
	return item, nil
}

// CheckIn checks file in, checkInType: 0 - Minor, 1 - Major, 2 - Overwrite
func (file *File) CheckIn(comment string, checkInType int) ([]byte, error) {
	endpoint := fmt.Sprintf(
		"%s/CheckIn(comment='%s',checkintype=%d)",
		file.endpoint,
		comment,
		checkInType,
	)
	client := NewHTTPClient(file.client)
	return client.Post(endpoint, nil, file.config)
}

// CheckOut checks file out
func (file *File) CheckOut() ([]byte, error) {
	endpoint := fmt.Sprintf("%s/CheckOut", file.endpoint)
	client := NewHTTPClient(file.client)
	return client.Post(endpoint, nil, file.config)
}

// UndoCheckOut undoes file check out
func (file *File) UndoCheckOut() ([]byte, error) {
	endpoint := fmt.Sprintf("%s/UndoCheckOut", file.endpoint)
	client := NewHTTPClient(file.client)
	return client.Post(endpoint, nil, file.config)
}

// Publish publishes a file
func (file *File) Publish(comment string) ([]byte, error) {
	endpoint := fmt.Sprintf(
		"%s/Publish(comment='%s')",
		file.endpoint,
		comment,
	)
	client := NewHTTPClient(file.client)
	return client.Post(endpoint, nil, file.config)
}

// UnPublish un-publishes a file
func (file *File) UnPublish(comment string) ([]byte, error) {
	endpoint := fmt.Sprintf(
		"%s/Publish(comment='%s')",
		file.endpoint,
		comment,
	)
	client := NewHTTPClient(file.client)
	return client.Post(endpoint, nil, file.config)
}

// GetReader gets file io.ReadCloser
func (file *File) GetReader() (io.ReadCloser, error) {
	siteURL := file.client.AuthCnfg.GetSiteURL()
	endpoint := fmt.Sprintf("%s/_api/web/%s/$value", siteURL, file.endpoint)

	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	// Apply context
	if file.config != nil && file.config.Context != nil {
		req = req.WithContext(file.config.Context)
	}

	req.TransferEncoding = []string{"null"}
	for key, value := range getConfHeaders(file.config) {
		req.Header.Set(key, value)
	}

	resp, err := file.client.Execute(req)
	if err != nil {
		return nil, err
	}
	return resp.Body, nil
}

// Download file bytes
func (file *File) Download() ([]byte, error) {
	body, err := file.GetReader()
	if err != nil {
		return nil, err
	}
	defer shut(body)

	data, err := io.ReadAll(body)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// MoveTo file to new location within the same site
func (file *File) MoveTo(newURL string, overwrite bool) ([]byte, error) {
	flag := 0
	if overwrite {
		flag = 1
	}
	endpoint := fmt.Sprintf("%s/MoveTo(newurl='%s',flags=%d)", file.endpoint, newURL, flag)
	client := NewHTTPClient(file.client)
	return client.Post(endpoint, nil, file.config)
}

// CopyTo file to new location within the same site
func (file *File) CopyTo(newURL string, overwrite bool) ([]byte, error) {
	endpoint := fmt.Sprintf("%s/CopyTo(strnewurl='%s',boverwrite=%t)", file.endpoint, newURL, overwrite)
	client := NewHTTPClient(file.client)
	return client.Post(endpoint, nil, file.config)
}

// ContextInfo ...
func (file *File) ContextInfo() (*ContextInfo, error) {
	return NewContext(file.client, file.ToURL(), file.config).Get()
}

// Props gets Properties API instance queryable collection for this File
func (file *File) Props() *Properties {
	return NewProperties(
		file.client,
		fmt.Sprintf("%s/Properties", file.endpoint),
		file.config,
		"file",
	)
}
