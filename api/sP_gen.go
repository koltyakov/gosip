// Package api :: This is auto generated file, do not edit manually
package api

// Conf receives custom request config definition, e.g. custom headers, custom OData mod
func (sP *SP) Conf(config *RequestConfig) *SP {
	sP.config = config
	return sP
}
