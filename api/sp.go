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

// NewSP ...
func NewSP(client *gosip.SPClient) *SP {
	return &SP{client: client}
}

// ToURL ...
func (sp *SP) ToURL() string {
	return sp.client.AuthCnfg.GetSiteURL()
}

// Conf ...
func (sp *SP) Conf(config *RequestConfig) *SP {
	sp.config = config
	return sp
}

// Web API object getter
func (sp *SP) Web() *Web {
	return NewWeb(
		sp.client,
		fmt.Sprintf("%s/_api/Web", sp.ToURL()),
		sp.config,
	)
}

// Site API object getter
func (sp *SP) Site() *Site {
	return NewSite(
		sp.client,
		fmt.Sprintf("%s/_api/Site", sp.ToURL()),
		sp.config,
	)
}

// Utility getter
func (sp *SP) Utility() *Utility {
	return NewUtility(sp.client, sp.ToURL(), sp.config)
}

// Search getter
func (sp *SP) Search() *Search {
	return NewSearch(
		sp.client,
		fmt.Sprintf("%s/_api/Search", sp.ToURL()),
		sp.config,
	)
}
