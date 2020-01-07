package api

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/koltyakov/gosip"
)

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

// Conf receives custom request config definition, e.g. custom headers, custom OData mod
func (file *File) Conf(config *RequestConfig) *File {
	file.config = config
	return file
}

// Select adds $select OData modifier
func (file *File) Select(oDataSelect string) *File {
	file.modifiers.AddSelect(oDataSelect)
	return file
}

// Expand adds $expand OData modifier
func (file *File) Expand(oDataExpand string) *File {
	file.modifiers.AddExpand(oDataExpand)
	return file
}

// Get gets file data object
func (file *File) Get() (FileResp, error) {
	sp := NewHTTPClient(file.client)
	return sp.Get(file.ToURL(), getConfHeaders(file.config))
}

// Update updates Field's metadata with properties provided in `body` parameter
// where `body` is byte array representation of JSON string payload relevalt to SP.File object
func (file *File) Update(body []byte) (FieldResp, error) {
	body = patchMetadataType(body, "SP.File")
	sp := NewHTTPClient(file.client)
	return sp.Update(file.endpoint, body, getConfHeaders(file.config))
}

// Delete deletes this file skipping recycle bin
func (file *File) Delete() error {
	sp := NewHTTPClient(file.client)
	_, err := sp.Delete(file.endpoint, getConfHeaders(file.config))
	return err
}

// Recycle moves this file to the recycle bin
func (file *File) Recycle() error {
	sp := NewHTTPClient(file.client)
	endpoint := fmt.Sprintf("%s/Recycle", file.endpoint)
	_, err := sp.Post(endpoint, nil, getConfHeaders(file.config))
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
	sp := NewHTTPClient(file.client)

	data, err := sp.Get(apiURL.String(), getConfHeaders(file.config))
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
	sp := NewHTTPClient(file.client)
	return sp.Post(endpoint, nil, getConfHeaders(file.config))
}

// CheckOut checks file out
func (file *File) CheckOut() ([]byte, error) {
	endpoint := fmt.Sprintf("%s/CheckOut", file.endpoint)
	sp := NewHTTPClient(file.client)
	return sp.Post(endpoint, nil, getConfHeaders(file.config))
}

// UndoCheckOut undoes file check out
func (file *File) UndoCheckOut() ([]byte, error) {
	endpoint := fmt.Sprintf("%s/UndoCheckOut", file.endpoint)
	sp := NewHTTPClient(file.client)
	return sp.Post(endpoint, nil, getConfHeaders(file.config))
}

// GetReader gets file io.ReadCloser
func (file *File) GetReader() (io.ReadCloser, error) {
	endpoint := fmt.Sprintf("%s/$value", file.endpoint)

	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
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
	defer body.Close()

	data, err := ioutil.ReadAll(body)
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
	sp := NewHTTPClient(file.client)
	return sp.Post(endpoint, nil, getConfHeaders(file.config))
}

// CopyTo file to new location within the same site
func (file *File) CopyTo(newURL string, overwrite bool) ([]byte, error) {
	endpoint := fmt.Sprintf("%s/CopyTo(strnewurl='%s',boverwrite=%t)", file.endpoint, newURL, overwrite)
	sp := NewHTTPClient(file.client)
	return sp.Post(endpoint, nil, getConfHeaders(file.config))
}

// ContextInfo ...
func (file *File) ContextInfo() (*ContextInfo, error) {
	return NewContext(file.client, file.ToURL(), file.config).Get()
}

/* Response helpers */

// Data : to get typed data
func (fileResp *FileResp) Data() *FileInfo {
	data := NormalizeODataItem(*fileResp)
	res := &FileInfo{}
	json.Unmarshal(data, &res)
	return res
}

// Unmarshal : to unmarshal to custom object
func (fileResp *FileResp) Unmarshal(obj interface{}) error {
	data := NormalizeODataItem(*fileResp)
	return json.Unmarshal(data, obj)
}
