// Package api :: This is auto generated file, do not edit manually
package api

// Conf receives custom request config definition, e.g. custom headers, custom OData mod
func (field *Field) Conf(config *RequestConfig) *Field {
	field.config = config
	return field
}

// Select adds $select OData modifier
func (field *Field) Select(oDataSelect string) *Field {
	field.modifiers.AddSelect(oDataSelect)
	return field
}

// Expand adds $expand OData modifier
func (field *Field) Expand(oDataExpand string) *Field {
	field.modifiers.AddExpand(oDataExpand)
	return field
}
