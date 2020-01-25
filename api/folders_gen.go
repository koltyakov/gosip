// Package api :: This is auto generated file, do not edit manually
package api

// Conf receives custom request config definition, e.g. custom headers, custom OData mod
func (folders *Folders) Conf(config *RequestConfig) *Folders {
	folders.config = config
	return folders
}

// Select adds $select OData modifier
func (folders *Folders) Select(oDataSelect string) *Folders {
	folders.modifiers.AddSelect(oDataSelect)
	return folders
}

// Expand adds $expand OData modifier
func (folders *Folders) Expand(oDataExpand string) *Folders {
	folders.modifiers.AddExpand(oDataExpand)
	return folders
}

// Filter adds $filter OData modifier
func (folders *Folders) Filter(oDataFilter string) *Folders {
	folders.modifiers.AddFilter(oDataFilter)
	return folders
}

// Top adds $top OData modifier
func (folders *Folders) Top(oDataTop int) *Folders {
	folders.modifiers.AddTop(oDataTop)
	return folders
}

// OrderBy adds $orderby OData modifier
func (folders *Folders) OrderBy(oDataOrderBy string, ascending bool) *Folders {
	folders.modifiers.AddOrderBy(oDataOrderBy, ascending)
	return folders
}
