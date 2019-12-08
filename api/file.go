package api

import (
	"encoding/json"
	"fmt"
	"net/url"
	"time"

	"github.com/koltyakov/gosip"
)

// File ...
type File struct {
	client    *gosip.SPClient
	config    *RequestConfig
	endpoint  string
	modifiers map[string]string
}

// FileInfo ...
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

// FileResp ...
type FileResp []byte

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
	query := apiURL.Query() // url.Values{}
	for k, v := range file.modifiers {
		query.Set(k, trimMultiline(v))
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
func (file *File) Get() (FileResp, error) {
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

// CheckIn file, checkInType: 0 - Minor, 1 - Major, 2 - Overwrite
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

// CheckOut file
func (file *File) CheckOut() ([]byte, error) {
	endpoint := fmt.Sprintf("%s/CheckOut", file.endpoint)
	sp := NewHTTPClient(file.client)
	return sp.Post(endpoint, nil, getConfHeaders(file.config))
}

// UndoCheckOut file
func (file *File) UndoCheckOut() ([]byte, error) {
	endpoint := fmt.Sprintf("%s/UndoCheckOut", file.endpoint)
	sp := NewHTTPClient(file.client)
	return sp.Post(endpoint, nil, getConfHeaders(file.config))
}

// ToDo:
// Move/Copy to
// Declare as record

/* Response helpers */

// Data : to get typed data
func (fileResp *FileResp) Data() *FileInfo {
	data := parseODataItem(*fileResp)
	res := &FileInfo{}
	json.Unmarshal(data, &res)
	return res
}

// Unmarshal : to unmarshal to custom object
func (fileResp *FileResp) Unmarshal(obj interface{}) error {
	data := parseODataItem(*fileResp)
	return json.Unmarshal(data, obj)
}
