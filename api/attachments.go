package api

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/koltyakov/gosip"
)

// Attachments represent SharePoint List Items Attachments API queryable collection struct
// Always use NewAttachments constructor instead of &Attachments{}
type Attachments struct {
	client   *gosip.SPClient
	config   *RequestConfig
	endpoint string
}

// Attachment represents SharePoint List Items Attachment API queryable object struct
// Always use NewAttachment constructor instead of &Attachment{}
type Attachment struct {
	client   *gosip.SPClient
	config   *RequestConfig
	endpoint string
}

// AttachmentInfo - attachment API response payload structure
type AttachmentInfo struct {
	FileName          string `json:"FileName"`
	ServerRelativeURL string `json:"ServerRelativeUrl"`
}

// AttachmentsResp - attachments response type with helper processor methods
type AttachmentsResp []byte

// AttachmentResp - attachment response type with helper processor methods
type AttachmentResp []byte

// NewAttachments - Attachments struct constructor function
func NewAttachments(client *gosip.SPClient, endpoint string, config *RequestConfig) *Attachments {
	return &Attachments{
		client:   client,
		endpoint: endpoint,
		config:   config,
	}
}

// NewAttachment - Attachment struct constructor function
func NewAttachment(client *gosip.SPClient, endpoint string, config *RequestConfig) *Attachment {
	return &Attachment{
		client:   client,
		endpoint: endpoint,
		config:   config,
	}
}

// Get gets attachments collection response
func (attachments *Attachments) Get() (AttachmentsResp, error) {
	sp := NewHTTPClient(attachments.client)
	return sp.Get(attachments.endpoint, getConfHeaders(attachments.config))
}

// Add uploads new attachment to the item
func (attachments *Attachments) Add(name string, content []byte) (AttachmentResp, error) {
	sp := NewHTTPClient(attachments.client)
	endpoint := fmt.Sprintf("%s/Add(FileName='%s')", attachments.endpoint, name)
	return sp.Post(endpoint, content, getConfHeaders(attachments.config))
}

// GetByName gets an attachment by its name
func (attachments *Attachments) GetByName(fileName string) *Attachment {
	return NewAttachment(
		attachments.client,
		fmt.Sprintf("%s('%s')", attachments.endpoint, fileName),
		attachments.config,
	)
}

// Get gets attachment data object
func (attachment *Attachment) Get() (AttachmentResp, error) {
	sp := NewHTTPClient(attachment.client)
	return sp.Get(attachment.endpoint, getConfHeaders(attachment.config))
}

// Delete delete an attachment skipping recycle bin
func (attachment *Attachment) Delete() error {
	sp := NewHTTPClient(attachment.client)
	_, err := sp.Delete(attachment.endpoint, getConfHeaders(attachment.config))
	return err
}

// Recycle moves an attachment to the recycle bin
func (attachment *Attachment) Recycle() error {
	sp := NewHTTPClient(attachment.client)
	endpoint := fmt.Sprintf("%s/RecycleObject", attachment.endpoint)
	_, err := sp.Post(endpoint, nil, getConfHeaders(attachment.config))
	return err
}

// GetReader gets attachment body data reader
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

// Dowload downloads attachment's as byte array
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
	collection, _ := normalizeODataCollection(*attachmentsResp)
	attachments := []AttachmentResp{}
	for _, attachment := range collection {
		attachments = append(attachments, AttachmentResp(attachment))
	}
	return attachments
}

// Normalized returns normalized body
func (attachmentsResp *AttachmentsResp) Normalized() []byte {
	normalized, _ := NormalizeODataCollection(*attachmentsResp)
	return normalized
}

// Data : to get typed data
func (attachmentResp *AttachmentResp) Data() *AttachmentInfo {
	data := NormalizeODataItem(*attachmentResp)
	res := &AttachmentInfo{}
	json.Unmarshal(data, &res)
	return res
}

// Normalized returns normalized body
func (attachmentResp *AttachmentResp) Normalized() []byte {
	return NormalizeODataItem(*attachmentResp)
}
