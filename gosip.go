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
	SetAuth(req *http.Request, client *SPClient) error // Authentication middleware fabric
	GetSiteURL() string                                // SiteURL getter method
	GetStrategy() string                               // Strategy code getter (triggered on demand)
	ReadConfig(configPath string) error                // Reads credentials from storage (triggered on demand)
	WriteConfig(configPath string) error               // Writes credential to storage (triggered on demand)
	GetAuth() (string, error)                          // Authentication initializer
}

// SPClient : SharePoint HTTP client struct
type SPClient struct {
	http.Client
	AuthCnfg   AuthCnfg // authentication configuration interface
	ConfigPath string   // private.json location path, optional when AuthCnfg is provided with creds explicitely

	RetryPolicies map[int]int // allows redefine error state requests retry policies
}

// RetryPolicies : error state requests default retry policies
var retryPolicies = map[int]int{
	401: 5, // on 401 :: Unauthorized
	429: 5, // on 429 :: Too many requests throttling error response
	500: 1, // on 500 :: Internal Server Error
	503: 5, // on 503 :: Service Unavailable
}

// Execute : SharePoint HTTP client
// is a wrapper for standard http.Client' `Do` method,
// injects authorization tokens, etc.
func (c *SPClient) Execute(req *http.Request) (*http.Response, error) {

	// Apply authentication flow
	if res, err := c.applyAuth(req); err != nil {
		return res, err
	}

	// Setup request default headers
	if err := c.applyHeaders(req); err != nil {
		// An error might occur only when calling for the digest
		res := &http.Response{
			Status:     "400 Bad Request",
			StatusCode: 400,
			Request:    req,
		}
		return res, err
	}

	// Sending actual request to SharePoint API/resource
	resp, err := c.Do(req)
	if err != nil {
		// Retry only for NTML
		if c.AuthCnfg.GetStrategy() == "ntlm" && c.shouldRetry(req, resp, 5) {
			return c.Execute(req)
		}
		return resp, err
	}

	// Wait and retry after a delay for error state responses, due to retry policies
	if retries := c.getRetryPolicy(resp.StatusCode); retries > 0 {
		// When it should, shouldRetry not only checks but waits before the retry
		if c.shouldRetry(req, resp, retries) {
			return c.Execute(req)
		}
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

// applyAuth applyes authentication flow
func (c *SPClient) applyAuth(req *http.Request) (*http.Response, error) {
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

	return nil, nil
}

// applyHeaders patches request readers for SP API defaults
func (c *SPClient) applyHeaders(req *http.Request) error {
	// Inject X-RequestDigest header when needed
	digestIsRequired := req.Method == "POST" &&
		!strings.Contains(strings.ToLower(req.URL.Path), "/_api/contextinfo") &&
		req.Header.Get("X-RequestDigest") == ""

	if digestIsRequired {
		digest, err := GetDigest(c)
		if err != nil {
			return err
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

	return nil
}

// getRetryPolicy receives retries policy retry number
func (c *SPClient) getRetryPolicy(statusCode int) int {
	// Apply default policies
	if c.RetryPolicies == nil {
		c.RetryPolicies = retryPolicies
	} else {
		// Append defaults to custom
		for status, retries := range retryPolicies {
			if _, ok := c.RetryPolicies[status]; !ok {
				c.RetryPolicies[status] = retries
			}
		}
	}
	return c.RetryPolicies[statusCode]
}

// shouldRetry checks should the request be retried, used with specific resp.StatusCode's
func (c *SPClient) shouldRetry(req *http.Request, resp *http.Response, retries int) bool {
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
