// Package api :: This is auto generated file, do not edit manually
package api

import "encoding/json"

/* Response helpers */

// Data response helper
func (recycleBinItemResp *RecycleBinItemResp) Data() *RecycleBinItemInfo {
	data := NormalizeODataItem(*recycleBinItemResp)
	res := &RecycleBinItemInfo{}
	json.Unmarshal(data, &res)
	return res
}

// Normalized returns normalized body
func (recycleBinItemResp *RecycleBinItemResp) Normalized() []byte {
	return NormalizeODataItem(*recycleBinItemResp)
}
