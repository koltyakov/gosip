/*
Package gosip is pure Go library for dealing with SharePoint unattended authentication and API consumption.

It supports a variety of different authentication strategies such as:
  - ADFS user credentials
  - Auth to SharePoint behind a reverse proxy (TMG, WAP)
  - Form-based authentication (FBA)
  - Addin only permissions
  - SAML based with user credentials

Amongst supported platform versions are:
  - SharePoint Online (SPO)
  - On-Prem: 2019, 2016, and 2013
*/
package gosip

import (
	"fmt"
	"net/http"
	"strings"
)

// AuthCnfg is an abstract auth config interface,
// allows different authentications strategies' dependency injection
type AuthCnfg interface {
	ReadConfig(configPath string) error
	WriteConfig(configPath string) error
	GetAuth() (string, error)
	GetSiteURL() string
	GetStrategy() string
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

	// Read stored creds and config
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

	// Wrap SharePoint authentication
	err := c.AuthCnfg.SetAuth(req, c)
	if err != nil {
		res := &http.Response{
			Status:     "401 Unauthorized",
			StatusCode: 401,
			Request:    req,
		}
		return res, err
	}

	// Inject X-RequestDigest header when needed
	digestIsRequired := req.Method == "POST" &&
		!strings.Contains(strings.ToLower(req.URL.Path), "/_api/contextinfo") &&
		req.Header.Get("X-RequestDigest") == ""
	if digestIsRequired {
		digest, err := GetDigest(c)
		if err != nil {
			res := &http.Response{
				Status:     "400 Bad Request",
				StatusCode: 400,
				Request:    req,
			}
			return res, err
		}
		req.Header.Add("X-RequestDigest", digest)
	}

	return c.Do(req)
}
