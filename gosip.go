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
	"io/ioutil"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const version = "1.0.0"

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

	// Can't resolve context siteURL
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
		req.Header.Set("X-RequestDigest", digest)
	}

	// Default SP REST API headers
	if req.Header.Get("Accept") == "" {
		req.Header.Set("Accept", "application/json")
	}
	if req.Header.Get("Content-Type") == "" {
		req.Header.Set("Content-Type", "application/json;odata=verbose;charset=utf-8")
	}

	// Vendor/client header
	if req.Header.Get("X-ClientService-ClientTag") == "" {
		req.Header.Set("X-ClientService-ClientTag", fmt.Sprintf("Gosip:@%s", version))
	}
	if req.Header.Get("User-Agent") == "" {
		req.Header.Set("User-Agent", fmt.Sprintf("NONISV|Go|Gosip/@%s", version))
	}

	resp, err := c.Do(req)
	if err != nil {
		if c.retry(req, resp, 5) {
			return c.Execute(req)
		}
		return resp, err
	}

	// Wait and retry after a delay on 429 :: Too many requests throttling error response
	if resp.StatusCode == 429 && c.retry(req, resp, 5) {
		return c.Execute(req)
	}

	// Wait and retry on 503 :: Service Unavailable
	if resp.StatusCode == 503 && c.retry(req, resp, 5) {
		return c.Execute(req)
	}

	// Wait and retry on 500 :: Internal Server Error
	if resp.StatusCode == 500 && c.retry(req, resp, 1) {
		return c.Execute(req)
	} // temporary workaround to fix unstable SPO service (https://github.com/SharePoint/sp-dev-docs/issues/4952)

	// Wait and retry on 401 :: Unauthorized
	if resp.StatusCode == 401 && c.retry(req, resp, 5) {
		return c.Execute(req)
	}

	// Return meaningful error message
	if err == nil && !(resp.StatusCode >= 200 && resp.StatusCode < 300) {
		defer resp.Body.Close()
		details, _ := ioutil.ReadAll(resp.Body)
		err = fmt.Errorf("%s :: %s", resp.Status, details)
		// Unescape unicode-escaped error messages for non Latin languages
		if unescaped, e := strconv.Unquote(`"` + strings.Replace(fmt.Sprintf("%s", details), `"`, `\"`, -1) + `"`); e == nil {
			err = fmt.Errorf("%s :: %s", resp.Status, unescaped)
		}
	}

	return resp, err
}

func (c *SPClient) retry(req *http.Request, resp *http.Response, retries int) bool {
	retry, _ := strconv.Atoi(req.Header.Get("X-Gosip-Retry"))
	if retry < retries {
		retryAfter := 0
		if resp != nil {
			retryAfter, _ = strconv.Atoi(resp.Header.Get("Retry-After"))
		}
		req.Header.Set("X-Gosip-Retry", strconv.Itoa(retry+1))
		if retryAfter != 0 {
			time.Sleep(time.Duration(retryAfter) * time.Second) // wait for Retry-After header info value
		} else {
			time.Sleep(time.Duration(100*math.Pow(2, float64(retry))) * time.Millisecond) // no Retry-After header
		}
		return true
	}
	return false
}
