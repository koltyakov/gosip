/*
Package gosip is pure Go library for dealing with SharePoint unattended authentication and API consumption.

It supports a variety of different authentication strategies such as:

* ADFS user credentials

* Auth to SharePoint behind a reverse proxy (TMG, WAP)

* Form-based authentication (FBA)

* Addin only permissions

* SAML based with user credentials


Amongst supported platform versions are:

* SharePoint Online (SPO)

* On-Prem: 2019, 2016, and 2013
*/
package gosip

import (
	"fmt"
	"net/http"
)

// AuthCnfg is an abstract auth config interface,
// allows different authentications strategies' dependency injection
type AuthCnfg interface {
	ReadConfig(configPath string) error
	GetAuth() (string, error)
	GetSiteURL() string
	SetAuth(req *http.Request, client *SPClient) error
}

// SPClient : SharePoint HTTP client struct
type SPClient struct {
	http.Client
	AuthCnfg   AuthCnfg
	ConfigPath string
}

// Execute : SharePoint HTTP client
// is a wrapper for standard http.Client' `Do` method,
// injects authorization tokens, etc.
func (c *SPClient) Execute(req *http.Request) (*http.Response, error) {
	if c.ConfigPath != "" && c.AuthCnfg.GetSiteURL() == "" {
		c.AuthCnfg.ReadConfig(c.ConfigPath)
	}
	if c.AuthCnfg.GetSiteURL() == "" {
		res := &http.Response{
			Status:     "400 Bad Request",
			StatusCode: 400,
			Request:    req,
		}
		return res, fmt.Errorf("client initialization error, no siteUrl is provided")
	}
	err := c.AuthCnfg.SetAuth(req, c)
	if err != nil {
		res := &http.Response{
			Status:     "401 Access Denied",
			StatusCode: 401,
			Request:    req,
		}
		return res, err
	}
	return c.Do(req)
}
