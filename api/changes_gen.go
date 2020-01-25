// Package api :: This is auto generated file, do not edit manually
package api

// Conf receives custom request config definition, e.g. custom headers, custom OData mod
func (changes *Changes) Conf(config *RequestConfig) *Changes {
	changes.config = config
	return changes
}
