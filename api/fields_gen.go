// Package api :: This is auto generated file, do not edit manually
package api

// Conf receives custom request config definition, e.g. custom headers, custom OData mod
func (fields *Fields) Conf(config *RequestConfig) *Fields {
	fields.config = config
	return fields
}

// Select adds $select OData modifier
func (fields *Fields) Select(oDataSelect string) *Fields {
	fields.modifiers.AddSelect(oDataSelect)
	return fields
}

// Expand adds $expand OData modifier
func (fields *Fields) Expand(oDataExpand string) *Fields {
	fields.modifiers.AddExpand(oDataExpand)
	return fields
}

// Filter adds $filter OData modifier
func (fields *Fields) Filter(oDataFilter string) *Fields {
	fields.modifiers.AddFilter(oDataFilter)
	return fields
}

// Top adds $top OData modifier
func (fields *Fields) Top(oDataTop int) *Fields {
	fields.modifiers.AddTop(oDataTop)
	return fields
}

// OrderBy adds $orderby OData modifier
func (fields *Fields) OrderBy(oDataOrderBy string, ascending bool) *Fields {
	fields.modifiers.AddOrderBy(oDataOrderBy, ascending)
	return fields
}

/* Response helpers */

// Data response helper
func (fieldsResp *FieldsResp) Data() []FieldResp {
	collection, _ := normalizeODataCollection(*fieldsResp)
	fields := []FieldResp{}
	for _, item := range collection {
		fields = append(fields, FieldResp(item))
	}
	return fields
}

// Normalized returns normalized body
func (fieldsResp *FieldsResp) Normalized() []byte {
	normalized, _ := NormalizeODataCollection(*fieldsResp)
	return normalized
}
