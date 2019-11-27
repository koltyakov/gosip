package api

import (
	"fmt"

	"github.com/koltyakov/gosip"
)

// SP - SharePoint REST API root struct
type SP struct {
	client *gosip.SPClient
	config *RequestConfig
}

// NewSPCtx ...
func NewSPCtx(client *gosip.SPClient) *SP {
	return &SP{client: client}
}

// Conf ...
func (sp *SP) Conf(config *RequestConfig) *SP {
	sp.config = config
	return sp
}

// Web API object getter
func (sp *SP) Web() *Web {
	return &Web{
		client:   sp.client,
		config:   sp.config,
		endpoint: fmt.Sprintf("%s/_api/Web", sp.client.AuthCnfg.GetSiteURL()),
	}
}
