/*
Package gosip is pure Go library for dealing with SharePoint unattended authentication and API consumption.

It supports a variety of different authentication strategies such as:
  - SAML based with user credentials
  - Add-in only permissions
  - Azure AD auth strategies (via extensions https://go.spflow.com/auth/custom-auth)
  - ADFS user credentials
  - NTLM/NTLM v2 windows auth
  - Auth to SharePoint behind a reverse proxy (TMG, WAP)
  - Form-based authentication (FBA)
  - Web login/On-Demand auth (via extension https://go.spflow.com/auth/custom-auth/on-demand)

Amongst supported platform versions are:
  - SharePoint Online (SPO)
  - On-Premise: 2019, 2016, and 2013
*/
package gosip

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const version = "1.0.0"

// AuthCnfg is an abstract auth config interface,
// allows different authentications strategies' dependency injection
type AuthCnfg interface {
	// GetAuth Authentication initializer (token/cookie/header, expiration, error)
	// to support capability for exposing tokens for external tools
	// e.g. as of this sample project https://github.com/koltyakov/spvault
	GetAuth() (string, int64, error)

	// SetAuth Authentication middleware fabric
	// applies round tripper or enriches requests with authentication and metadata
	SetAuth(req *http.Request, client *SPClient) error

	ParseConfig(jsonConf []byte) error  // Parses credentials from a provided JSON byte array content
	ReadConfig(configPath string) error // Reads credentials from storage

	GetSiteURL() string  // SiteURL getter method
	GetStrategy() string // Strategy code getter

	// WriteConfig(configPath string) error // Writes credential to storage // considering remove the method from interface
}

// SPClient : SharePoint HTTP client struct
type SPClient struct {
	http.Client
	AuthCnfg   AuthCnfg // authentication configuration interface
	ConfigPath string   // private.json location path, optional when AuthCnfg is provided with creds explicitly

	RetryPolicies map[int]int   // allows redefining error state requests retry policies
	Hooks         *HookHandlers // hook handlers definition
}

// Execute : SharePoint HTTP client
// is a wrapper for standard http.Client' `Do` method, injects authorization tokens, etc.
func (c *SPClient) Execute(req *http.Request) (*http.Response, error) {
	reqTime := time.Now()

	// Apply authentication flow
	if res, err := c.applyAuth(req); err != nil {
		c.onError(req, reqTime, 0, err)
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
		c.onError(req, reqTime, 0, err)
		return res, err
	}

	c.onRequest(req, reqTime, 0, nil)
	reqTime = time.Now() // update request time to exclude auth-related timings

	// Creating backup reader to be able to retry none nil body requests
	var bodyBackup io.Reader
	if req.Body != nil {
		var buf bytes.Buffer
		tee := io.TeeReader(req.Body, &buf)
		bodyBackup = &buf
		req.Body = io.NopCloser(tee)
	}

	// Sending actual request to SharePoint API/resource
	resp, err := c.Do(req)
	if err != nil {
		// Retry only for NTLM
		if c.AuthCnfg.GetStrategy() == "ntlm" && c.shouldRetry(req, resp, 5) {
			statusCode := 400
			if resp != nil {
				statusCode = resp.StatusCode
			}
			c.onRetry(req, reqTime, statusCode, nil)
			// Reset body reader closer
			if bodyBackup != nil {
				req.Body = io.NopCloser(bodyBackup)
			}
			return c.Execute(req)
		}
		c.onError(req, reqTime, 0, err)
		return resp, err
	}

	// Wait and retry after a delay for error state responses, due to retry policies
	if retries := c.getRetryPolicy(resp.StatusCode); retries > 0 {
		// Register retry in OnError hook
		// otherwise it only called in OnRetry after timeout right before the next call
		if resp.StatusCode == 429 {
			noRetry := req.Header.Get("X-Gosip-NoRetry")
			retry, _ := strconv.Atoi(req.Header.Get("X-Gosip-Retry"))
			if retry < retries && noRetry != "true" {
				c.onError(req, reqTime, resp.StatusCode, nil)
			}
		}

		// When it should, shouldRetry not only checks but waits before a retry
		if c.shouldRetry(req, resp, retries) {
			c.onRetry(req, reqTime, resp.StatusCode, nil)
			// Reset body reader closer
			if bodyBackup != nil {
				req.Body = io.NopCloser(bodyBackup)
			}
			return c.Execute(req)
		}
	}

	// Return meaningful error message
	if !(resp.StatusCode >= 200 && resp.StatusCode < 300) {
		var buf bytes.Buffer
		tee := io.TeeReader(resp.Body, &buf)
		details, _ := io.ReadAll(tee)
		err = fmt.Errorf("%s :: %s", resp.Status, details)
		// Unescape unicode-escaped error messages for non Latin languages
		if unescaped, e := strconv.Unquote(`"` + strings.Replace(string(details), `"`, `\"`, -1) + `"`); e == nil {
			err = fmt.Errorf("%s :: %s", resp.Status, unescaped)
		}
		resp.Body = io.NopCloser(&buf)
		c.onError(req, reqTime, resp.StatusCode, err)
	}

	c.onResponse(req, reqTime, resp.StatusCode, err)
	return resp, err
}

// applyAuth applies authentication flow
func (c *SPClient) applyAuth(req *http.Request) (*http.Response, error) {
	// Read stored credentials and config
	if c.ConfigPath != "" && c.AuthCnfg.GetSiteURL() == "" {
		_ = c.AuthCnfg.ReadConfig(c.ConfigPath)
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
	digestIsRequired := (req.Method == "POST" || req.Method == "PATCH" || req.Method == "MERGE") &&
		!strings.Contains(strings.ToLower(req.URL.Path), "/_api/contextinfo") &&
		req.Header.Get("X-RequestDigest") == ""

	if digestIsRequired {
		digest, err := GetDigest(req.Context(), c)
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
