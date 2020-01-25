// Package api :: This is auto generated file, do not edit manually
package api

// Conf receives custom request config definition, e.g. custom headers, custom OData mod
func (site *Site) Conf(config *RequestConfig) *Site {
	site.config = config
	return site
}

// Select adds $select OData modifier
func (site *Site) Select(oDataSelect string) *Site {
	site.modifiers.AddSelect(oDataSelect)
	return site
}

// Expand adds $expand OData modifier
func (site *Site) Expand(oDataExpand string) *Site {
	site.modifiers.AddExpand(oDataExpand)
	return site
}
