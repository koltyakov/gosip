// Package api :: This is auto generated file, do not edit manually
package api

import "encoding/json"

/* Response helpers */

// Data response helper
func (fieldLinkResp *FieldLinkResp) Data() *FieldLinkInfo {
	data := NormalizeODataItem(*fieldLinkResp)
	res := &FieldLinkInfo{}
	json.Unmarshal(data, &res)
	return res
}

// Normalized returns normalized body
func (fieldLinkResp *FieldLinkResp) Normalized() []byte {
	return NormalizeODataItem(*fieldLinkResp)
}
