// Package api :: This is auto generated file, do not edit manually
package api

// Conf receives custom request config definition, e.g. custom headers, custom OData mod
func (customActions *CustomActions) Conf(config *RequestConfig) *CustomActions {
	customActions.config = config
	return customActions
}

// Select adds $select OData modifier
func (customActions *CustomActions) Select(oDataSelect string) *CustomActions {
	customActions.modifiers.AddSelect(oDataSelect)
	return customActions
}

// Filter adds $filter OData modifier
func (customActions *CustomActions) Filter(oDataFilter string) *CustomActions {
	customActions.modifiers.AddFilter(oDataFilter)
	return customActions
}

// Top adds $top OData modifier
func (customActions *CustomActions) Top(oDataTop int) *CustomActions {
	customActions.modifiers.AddTop(oDataTop)
	return customActions
}

// OrderBy adds $orderby OData modifier
func (customActions *CustomActions) OrderBy(oDataOrderBy string, ascending bool) *CustomActions {
	customActions.modifiers.AddOrderBy(oDataOrderBy, ascending)
	return customActions
}
