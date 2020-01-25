// Package api :: This is auto generated file, do not edit manually
package api

// Conf receives custom request config definition, e.g. custom headers, custom OData mod
func (contentTypes *ContentTypes) Conf(config *RequestConfig) *ContentTypes {
	contentTypes.config = config
	return contentTypes
}

// Select adds $select OData modifier
func (contentTypes *ContentTypes) Select(oDataSelect string) *ContentTypes {
	contentTypes.modifiers.AddSelect(oDataSelect)
	return contentTypes
}

// Expand adds $expand OData modifier
func (contentTypes *ContentTypes) Expand(oDataExpand string) *ContentTypes {
	contentTypes.modifiers.AddExpand(oDataExpand)
	return contentTypes
}

// Filter adds $filter OData modifier
func (contentTypes *ContentTypes) Filter(oDataFilter string) *ContentTypes {
	contentTypes.modifiers.AddFilter(oDataFilter)
	return contentTypes
}

// Top adds $top OData modifier
func (contentTypes *ContentTypes) Top(oDataTop int) *ContentTypes {
	contentTypes.modifiers.AddTop(oDataTop)
	return contentTypes
}

// OrderBy adds $orderby OData modifier
func (contentTypes *ContentTypes) OrderBy(oDataOrderBy string, ascending bool) *ContentTypes {
	contentTypes.modifiers.AddOrderBy(oDataOrderBy, ascending)
	return contentTypes
}
