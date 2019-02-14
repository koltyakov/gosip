package gosip

import (
	"fmt"
	"net/http"
)

// AuthCnfg : abstract auth interface
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

// Execute : SharePoint HTTP client Do method
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
