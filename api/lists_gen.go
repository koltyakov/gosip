// Package api :: This is auto generated file, do not edit manually
package api

// Conf receives custom request config definition, e.g. custom headers, custom OData mod
func (lists *Lists) Conf(config *RequestConfig) *Lists {
	lists.config = config
	return lists
}

// Select adds $select OData modifier
func (lists *Lists) Select(oDataSelect string) *Lists {
	lists.modifiers.AddSelect(oDataSelect)
	return lists
}

// Expand adds $expand OData modifier
func (lists *Lists) Expand(oDataExpand string) *Lists {
	lists.modifiers.AddExpand(oDataExpand)
	return lists
}

// Filter adds $filter OData modifier
func (lists *Lists) Filter(oDataFilter string) *Lists {
	lists.modifiers.AddFilter(oDataFilter)
	return lists
}

// Top adds $top OData modifier
func (lists *Lists) Top(oDataTop int) *Lists {
	lists.modifiers.AddTop(oDataTop)
	return lists
}

// OrderBy adds $orderby OData modifier
func (lists *Lists) OrderBy(oDataOrderBy string, ascending bool) *Lists {
	lists.modifiers.AddOrderBy(oDataOrderBy, ascending)
	return lists
}

/* Response helpers */

// Data response helper
func (listsResp *ListsResp) Data() []ListResp {
	collection, _ := normalizeODataCollection(*listsResp)
	lists := []ListResp{}
	for _, item := range collection {
		lists = append(lists, ListResp(item))
	}
	return lists
}

// Normalized returns normalized body
func (listsResp *ListsResp) Normalized() []byte {
	normalized, _ := NormalizeODataCollection(*listsResp)
	return normalized
}
