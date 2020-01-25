// Package api :: This is auto generated file, do not edit manually
package api

// Conf receives custom request config definition, e.g. custom headers, custom OData mod
func (users *Users) Conf(config *RequestConfig) *Users {
	users.config = config
	return users
}

// Select adds $select OData modifier
func (users *Users) Select(oDataSelect string) *Users {
	users.modifiers.AddSelect(oDataSelect)
	return users
}

// Expand adds $expand OData modifier
func (users *Users) Expand(oDataExpand string) *Users {
	users.modifiers.AddExpand(oDataExpand)
	return users
}

// Filter adds $filter OData modifier
func (users *Users) Filter(oDataFilter string) *Users {
	users.modifiers.AddFilter(oDataFilter)
	return users
}

// Top adds $top OData modifier
func (users *Users) Top(oDataTop int) *Users {
	users.modifiers.AddTop(oDataTop)
	return users
}

// OrderBy adds $orderby OData modifier
func (users *Users) OrderBy(oDataOrderBy string, ascending bool) *Users {
	users.modifiers.AddOrderBy(oDataOrderBy, ascending)
	return users
}
