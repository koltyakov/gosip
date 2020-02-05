package api

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/koltyakov/gosip"
)

//go:generate ggen -ent Attachments -item Attachment -coll -helpers Data,Normalized
//go:generate ggen -ent Attachment -helpers Data,Normalized

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
	client := NewHTTPClient(attachments.client)
	return client.Get(attachments.endpoint, attachments.config)
}

// Add uploads new attachment to the item
func (attachments *Attachments) Add(name string, content io.Reader) (AttachmentResp, error) {
	client := NewHTTPClient(attachments.client)
	endpoint := fmt.Sprintf("%s/Add(FileName='%s')", attachments.endpoint, name)
	return client.Post(endpoint, content, attachments.config)
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
	client := NewHTTPClient(attachment.client)
	return client.Get(attachment.endpoint, attachment.config)
}

// Delete delete an attachment skipping recycle bin
func (attachment *Attachment) Delete() error {
	client := NewHTTPClient(attachment.client)
	_, err := client.Delete(attachment.endpoint, attachment.config)
	return err
}

// Recycle moves an attachment to the recycle bin
func (attachment *Attachment) Recycle() error {
	client := NewHTTPClient(attachment.client)
	endpoint := fmt.Sprintf("%s/RecycleObject", attachment.endpoint)
	_, err := client.Post(endpoint, nil, attachment.config)
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
