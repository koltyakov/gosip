// Package api :: This is auto generated file, do not edit manually
package api

// Conf receives custom request config definition, e.g. custom headers, custom OData mod
func (view *View) Conf(config *RequestConfig) *View {
	view.config = config
	return view
}

// Select adds $select OData modifier
func (view *View) Select(oDataSelect string) *View {
	view.modifiers.AddSelect(oDataSelect)
	return view
}

// Expand adds $expand OData modifier
func (view *View) Expand(oDataExpand string) *View {
	view.modifiers.AddExpand(oDataExpand)
	return view
}
