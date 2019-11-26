package api

import (
	"net/url"

	"github.com/koltyakov/gosip"
)

// Web ...
type Web struct {
	client   *gosip.SPClient
	conf     *Conf
	endpoint string
	oSelect  string
	oExpand  string
}

// Conf ...
func (web *Web) Conf(conf *Conf) *Web {
	web.conf = conf
	return web
}

// Select ...
func (web *Web) Select(oDataSelect string) *Web {
	web.oSelect = oDataSelect
	return web
}

// Expand ...
func (web *Web) Expand(oDataExpand string) *Web {
	web.oExpand = oDataExpand
	return web
}

// Get ...
func (web *Web) Get() ([]byte, error) {
	apiURL, _ := url.Parse(web.endpoint)
	query := url.Values{}
	if web.oSelect != "" {
		query.Add("$select", TrimMultiline(web.oSelect))
	}
	if web.oExpand != "" {
		query.Add("$expand", TrimMultiline(web.oExpand))
	}
	apiURL.RawQuery = query.Encode()
	sp := &HTTPClient{SPClient: web.client}
	return sp.Get(apiURL.String(), GetConfHeaders(web.conf))
}

// Delete ...
func (web *Web) Delete() ([]byte, error) {
	sp := &HTTPClient{SPClient: web.client}
	return sp.Delete(web.endpoint, GetConfHeaders(web.conf))
}

// Update ...
func (web *Web) Update(body []byte) ([]byte, error) {
	sp := &HTTPClient{SPClient: web.client}
	return sp.Update(web.endpoint, body, GetConfHeaders(web.conf))
}
