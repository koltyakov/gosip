// Package api :: This is auto generated file, do not edit manually
package api

// Conf receives custom request config definition, e.g. custom headers, custom OData mod
func (views *Views) Conf(config *RequestConfig) *Views {
	views.config = config
	return views
}

// Select adds $select OData modifier
func (views *Views) Select(oDataSelect string) *Views {
	views.modifiers.AddSelect(oDataSelect)
	return views
}

// Expand adds $expand OData modifier
func (views *Views) Expand(oDataExpand string) *Views {
	views.modifiers.AddExpand(oDataExpand)
	return views
}

// Filter adds $filter OData modifier
func (views *Views) Filter(oDataFilter string) *Views {
	views.modifiers.AddFilter(oDataFilter)
	return views
}

// Top adds $top OData modifier
func (views *Views) Top(oDataTop int) *Views {
	views.modifiers.AddTop(oDataTop)
	return views
}

// OrderBy adds $orderby OData modifier
func (views *Views) OrderBy(oDataOrderBy string, ascending bool) *Views {
	views.modifiers.AddOrderBy(oDataOrderBy, ascending)
	return views
}
