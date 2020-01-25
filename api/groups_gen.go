// Package api :: This is auto generated file, do not edit manually
package api

// Conf receives custom request config definition, e.g. custom headers, custom OData mod
func (groups *Groups) Conf(config *RequestConfig) *Groups {
	groups.config = config
	return groups
}

// Select adds $select OData modifier
func (groups *Groups) Select(oDataSelect string) *Groups {
	groups.modifiers.AddSelect(oDataSelect)
	return groups
}

// Expand adds $expand OData modifier
func (groups *Groups) Expand(oDataExpand string) *Groups {
	groups.modifiers.AddExpand(oDataExpand)
	return groups
}

// Filter adds $filter OData modifier
func (groups *Groups) Filter(oDataFilter string) *Groups {
	groups.modifiers.AddFilter(oDataFilter)
	return groups
}

// Top adds $top OData modifier
func (groups *Groups) Top(oDataTop int) *Groups {
	groups.modifiers.AddTop(oDataTop)
	return groups
}

// OrderBy adds $orderby OData modifier
func (groups *Groups) OrderBy(oDataOrderBy string, ascending bool) *Groups {
	groups.modifiers.AddOrderBy(oDataOrderBy, ascending)
	return groups
}
