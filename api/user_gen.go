// Package api :: This is auto generated file, do not edit manually
package api

import "encoding/json"

// Conf receives custom request config definition, e.g. custom headers, custom OData mod
func (user *User) Conf(config *RequestConfig) *User {
	user.config = config
	return user
}

// Select adds $select OData modifier
func (user *User) Select(oDataSelect string) *User {
	user.modifiers.AddSelect(oDataSelect)
	return user
}

// Expand adds $expand OData modifier
func (user *User) Expand(oDataExpand string) *User {
	user.modifiers.AddExpand(oDataExpand)
	return user
}

/* Response helpers */

// Data response helper
func (userResp *UserResp) Data() *UserInfo {
	data := NormalizeODataItem(*userResp)
	res := &UserInfo{}
	json.Unmarshal(data, &res)
	return res
}

// Normalized returns normalized body
func (userResp *UserResp) Normalized() []byte {
	return NormalizeODataItem(*userResp)
}
