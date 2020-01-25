// Package api :: This is auto generated file, do not edit manually
package api

// Conf receives custom request config definition, e.g. custom headers, custom OData mod
func (folder *Folder) Conf(config *RequestConfig) *Folder {
	folder.config = config
	return folder
}

// Select adds $select OData modifier
func (folder *Folder) Select(oDataSelect string) *Folder {
	folder.modifiers.AddSelect(oDataSelect)
	return folder
}

// Expand adds $expand OData modifier
func (folder *Folder) Expand(oDataExpand string) *Folder {
	folder.modifiers.AddExpand(oDataExpand)
	return folder
}
