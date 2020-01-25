// Package api :: This is auto generated file, do not edit manually
package api

import "encoding/json"

/* Response helpers */

// Data response helper
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
