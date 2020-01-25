// Package api :: This is auto generated file, do not edit manually
package api

// Conf receives custom request config definition, e.g. custom headers, custom OData mod
func (item *Item) Conf(config *RequestConfig) *Item {
	item.config = config
	return item
}

// Select adds $select OData modifier
func (item *Item) Select(oDataSelect string) *Item {
	item.modifiers.AddSelect(oDataSelect)
	return item
}

// Expand adds $expand OData modifier
func (item *Item) Expand(oDataExpand string) *Item {
	item.modifiers.AddExpand(oDataExpand)
	return item
}

/* Response helpers */

// Normalized returns normalized body
func (itemResp *ItemResp) Normalized() []byte {
	return NormalizeODataItem(*itemResp)
}
