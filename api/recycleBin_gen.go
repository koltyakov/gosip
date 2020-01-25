// Package api :: This is auto generated file, do not edit manually
package api

// Conf receives custom request config definition, e.g. custom headers, custom OData mod
func (recycleBin *RecycleBin) Conf(config *RequestConfig) *RecycleBin {
	recycleBin.config = config
	return recycleBin
}

// Select adds $select OData modifier
func (recycleBin *RecycleBin) Select(oDataSelect string) *RecycleBin {
	recycleBin.modifiers.AddSelect(oDataSelect)
	return recycleBin
}

// Expand adds $expand OData modifier
func (recycleBin *RecycleBin) Expand(oDataExpand string) *RecycleBin {
	recycleBin.modifiers.AddExpand(oDataExpand)
	return recycleBin
}

// Filter adds $filter OData modifier
func (recycleBin *RecycleBin) Filter(oDataFilter string) *RecycleBin {
	recycleBin.modifiers.AddFilter(oDataFilter)
	return recycleBin
}

// Top adds $top OData modifier
func (recycleBin *RecycleBin) Top(oDataTop int) *RecycleBin {
	recycleBin.modifiers.AddTop(oDataTop)
	return recycleBin
}

// OrderBy adds $orderby OData modifier
func (recycleBin *RecycleBin) OrderBy(oDataOrderBy string, ascending bool) *RecycleBin {
	recycleBin.modifiers.AddOrderBy(oDataOrderBy, ascending)
	return recycleBin
}

/* Response helpers */

// Data response helper
func (recycleBinResp *RecycleBinResp) Data() []RecycleBinItemResp {
	collection, _ := normalizeODataCollection(*recycleBinResp)
	recycleBin := []RecycleBinItemResp{}
	for _, item := range collection {
		recycleBin = append(recycleBin, RecycleBinItemResp(item))
	}
	return recycleBin
}

// Normalized returns normalized body
func (recycleBinResp *RecycleBinResp) Normalized() []byte {
	normalized, _ := NormalizeODataCollection(*recycleBinResp)
	return normalized
}
