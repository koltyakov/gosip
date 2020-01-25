// Package api :: This is auto generated file, do not edit manually
package api

// Conf receives custom request config definition, e.g. custom headers, custom OData mod
func (properties *Properties) Conf(config *RequestConfig) *Properties {
	properties.config = config
	return properties
}

// Select adds $select OData modifier
func (properties *Properties) Select(oDataSelect string) *Properties {
	properties.modifiers.AddSelect(oDataSelect)
	return properties
}

// Expand adds $expand OData modifier
func (properties *Properties) Expand(oDataExpand string) *Properties {
	properties.modifiers.AddExpand(oDataExpand)
	return properties
}
