package api

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/koltyakov/gosip"
)

// Attachments ...
type Attachments struct {
	client   *gosip.SPClient
	config   *RequestConfig
	endpoint string
}

// Attachment ...
type Attachment struct {
	client   *gosip.SPClient
	config   *RequestConfig
	endpoint string
}

// AttachmentInfo ...
type AttachmentInfo struct {
	FileName          string `json:"FileName"`
	ServerRelativeURL string `json:"ServerRelativeUrl"`
}

// AttachmentsResp ...
type AttachmentsResp []byte

// AttachmentResp ...
type AttachmentResp []byte

// NewAttachments ...
func NewAttachments(client *gosip.SPClient, endpoint string, config *RequestConfig) *Attachments {
	return &Attachments{
		client:   client,
		endpoint: endpoint,
		config:   config,
	}
}

// NewAttachment ...
func NewAttachment(client *gosip.SPClient, endpoint string, config *RequestConfig) *Attachment {
	return &Attachment{
		client:   client,
		endpoint: endpoint,
		config:   config,
	}
}

// Get ...
func (attachments *Attachments) Get() (AttachmentsResp, error) {
	sp := NewHTTPClient(attachments.client)
	return sp.Get(attachments.endpoint, getConfHeaders(attachments.config))
}

// Add ...
func (attachments *Attachments) Add(name string, content []byte) (FileResp, error) {
	sp := NewHTTPClient(attachments.client)
	endpoint := fmt.Sprintf("%s/Add(FileName='%s')", attachments.endpoint, name)
	return sp.Post(endpoint, content, getConfHeaders(attachments.config))
}

// GetByName ...
func (attachments *Attachments) GetByName(fileName string) *Attachment {
	return NewAttachment(
		attachments.client,
		fmt.Sprintf("%s('%s')", attachments.endpoint, fileName),
		attachments.config,
	)
}

// Get ...
func (attachment *Attachment) Get() (AttachmentResp, error) {
	sp := NewHTTPClient(attachment.client)
	return sp.Get(attachment.endpoint, getConfHeaders(attachment.config))
}

// Delete ...
func (attachment *Attachment) Delete() ([]byte, error) {
	sp := NewHTTPClient(attachment.client)
	return sp.Delete(attachment.endpoint, getConfHeaders(attachment.config))
}

// Recycle ...
func (attachment *Attachment) Recycle() ([]byte, error) {
	sp := NewHTTPClient(attachment.client)
	endpoint := fmt.Sprintf("%s/RecycleObject", attachment.endpoint)
	return sp.Post(endpoint, nil, getConfHeaders(attachment.config))
}

// GetReader ...
func (attachment *Attachment) GetReader() (io.ReadCloser, error) {
	endpoint := fmt.Sprintf("%s/$value", attachment.endpoint)

	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	req.TransferEncoding = []string{"null"}
	for key, value := range getConfHeaders(attachment.config) {
		req.Header.Set(key, value)
	}

	resp, err := attachment.client.Execute(req)
	if err != nil {
		return nil, err
	}
	return resp.Body, err
}

// Dowload ...
func (attachment *Attachment) Dowload() ([]byte, error) {
	body, err := attachment.GetReader()
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

/* Response helpers */

// Data : to get typed data
func (attachmentsResp *AttachmentsResp) Data() []AttachmentResp {
	collection, _ := parseODataCollection(*attachmentsResp)
	attachments := []AttachmentResp{}
	for _, attachment := range collection {
		attachments = append(attachments, AttachmentResp(attachment))
	}
	return attachments
}

// Unmarshal : to unmarshal to custom object
func (attachmentsResp *AttachmentsResp) Unmarshal(obj interface{}) error {
	// collection := parseODataCollection(*ctsResp)
	// data, _ := json.Marshal(collection)
	data, _ := parseODataCollectionPlain(*attachmentsResp)
	return json.Unmarshal(data, obj)
}

// Data : to get typed data
func (attachmentResp *AttachmentResp) Data() *AttachmentInfo {
	data := parseODataItem(*attachmentResp)
	res := &AttachmentInfo{}
	json.Unmarshal(data, &res)
	return res
}

// Unmarshal : to unmarshal to custom object
func (attachmentResp *AttachmentResp) Unmarshal(obj interface{}) error {
	data := parseODataItem(*attachmentResp)
	return json.Unmarshal(data, obj)
}
