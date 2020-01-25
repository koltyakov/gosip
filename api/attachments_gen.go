// Package api :: This is auto generated file, do not edit manually
package api

/* Response helpers */

// Data response helper
func (attachmentsResp *AttachmentsResp) Data() []AttachmentResp {
	collection, _ := normalizeODataCollection(*attachmentsResp)
	attachments := []AttachmentResp{}
	for _, item := range collection {
		attachments = append(attachments, AttachmentResp(item))
	}
	return attachments
}

// Normalized returns normalized body
func (attachmentsResp *AttachmentsResp) Normalized() []byte {
	normalized, _ := NormalizeODataCollection(*attachmentsResp)
	return normalized
}
