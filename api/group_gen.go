// Package api :: This is auto generated file, do not edit manually
package api

import "encoding/json"

// Conf receives custom request config definition, e.g. custom headers, custom OData mod
func (group *Group) Conf(config *RequestConfig) *Group {
	group.config = config
	return group
}

// Select adds $select OData modifier
func (group *Group) Select(oDataSelect string) *Group {
	group.modifiers.AddSelect(oDataSelect)
	return group
}

// Expand adds $expand OData modifier
func (group *Group) Expand(oDataExpand string) *Group {
	group.modifiers.AddExpand(oDataExpand)
	return group
}

/* Response helpers */

// Data response helper
func (groupResp *GroupResp) Data() *GroupInfo {
	data := NormalizeODataItem(*groupResp)
	res := &GroupInfo{}
	json.Unmarshal(data, &res)
	return res
}

// Normalized returns normalized body
func (groupResp *GroupResp) Normalized() []byte {
	return NormalizeODataItem(*groupResp)
}
