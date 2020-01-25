// Package api :: This is auto generated file, do not edit manually
package api

// Conf receives custom request config definition, e.g. custom headers, custom OData mod
func (web *Web) Conf(config *RequestConfig) *Web {
	web.config = config
	return web
}

// Select adds $select OData modifier
func (web *Web) Select(oDataSelect string) *Web {
	web.modifiers.AddSelect(oDataSelect)
	return web
}

// Expand adds $expand OData modifier
func (web *Web) Expand(oDataExpand string) *Web {
	web.modifiers.AddExpand(oDataExpand)
	return web
}
