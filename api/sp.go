package api

import (
	"fmt"

	"github.com/koltyakov/gosip"
)

// SP - SharePoint REST API root struct
type SP struct {
	SPClient *gosip.SPClient
	conf     *Conf
}

// Conf struct
type Conf struct {
	Headers map[string]string
}

// Conf ...
func (sp *SP) Conf(conf *Conf) *SP {
	sp.conf = conf
	return sp
}

// Web API object getter
func (sp *SP) Web() *Web {
	return &Web{
		client:   sp.SPClient,
		conf:     sp.conf,
		endpoint: fmt.Sprintf("%s/_api/Web", sp.SPClient.AuthCnfg.GetSiteURL()),
	}
}
