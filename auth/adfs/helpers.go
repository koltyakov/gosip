package adfs

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/patrickmn/go-cache"

	"github.com/koltyakov/gosip/templates"
)

var (
	storage = cache.New(5*time.Minute, 10*time.Minute)
)

// GetAuth : get auth
func GetAuth(c *AuthCnfg) (string, error) {
	parsedURL, err := url.Parse(c.SiteURL)
	if err != nil {
		return "", err
	}

	cacheKey := parsedURL.Host + "@adfs@" + c.Username + "@" + c.Password
	if authCookie, found := storage.Get(cacheKey); found {
		return authCookie.(string), nil
	}

	var authCookie, expires string
	var expiry time.Duration

	// In case of WAP
	if c.AdfsCookie == "EdgeAccessCookie" {
		authCookie, expires, err = wapAuthFlow(c)
		if err != nil {
			return "", err
		}
		if expires == "" {
			expiry = 30 * time.Minute // ToDO: move to settings or dynamically get
		}
	} else {
		authCookie, expires, err = adfsAuthFlow(c, "")
		if err != nil {
			return "", err
		}
		expiresTime, _ := time.Parse(time.RFC3339, expires)
		expiry = time.Until(expiresTime) - 60*time.Second
	}

	storage.Set(cacheKey, authCookie, expiry)

	return authCookie, nil
}

func adfsAuthFlow(c *AuthCnfg, edgeCookie string) (string, string, error) {
	parsedAdfsURL, err := url.Parse(c.AdfsURL)
	if err != nil {
		return "", "", err
	}

	usernameMixedURL := fmt.Sprintf("%s://%s/adfs/services/trust/13/usernamemixed", parsedAdfsURL.Scheme, parsedAdfsURL.Host)
	samlBody, err := templates.AdfsSamlWsfedTemplate(usernameMixedURL, c.Username, c.Password, c.RelyingParty)
	if err != nil {
		return "", "", err
	}

	req, err := http.NewRequest("POST", usernameMixedURL, bytes.NewBuffer([]byte(samlBody)))
	if err != nil {
		return "", "", err
	}

	req.Header.Set("Content-Type", "application/soap+xml;charset=utf-8")
	if edgeCookie != "" {
		req.Header.Set("Cookie", edgeCookie)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", "", err
	}
	defer func() { _ = resp.Body.Close() }()

	res, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", "", err
	}

	// fmt.Printf("ADFS: %s\n", string(res))

	type samlAssertion struct {
		Fault    string `xml:"Body>Fault>Reason>Text"`
		Response struct {
			Token struct {
				Inner      []byte `xml:",innerxml"`
				Conditions struct {
					NotBefore    string `xml:"NotBefore,attr"`
					NotOnOrAfter string `xml:"NotOnOrAfter,attr"`
				} `xml:"Assertion>Conditions"`
			} `xml:"RequestedSecurityToken"`
			// This is for WPA (urn:AppProxy:com)
			Lifetime struct {
				Created string `xml:"Created"`
				Expires string `xml:"Expires"`
			} `xml:"Lifetime"`
		} `xml:"Body>RequestSecurityTokenResponseCollection>RequestSecurityTokenResponse"`
	}
	result := &samlAssertion{}
	if err := xml.Unmarshal(res, &result); err != nil {
		return "", "", err
	}

	if result.Fault != "" {
		return "", "", errors.New(result.Fault)
	}

	created := result.Response.Token.Conditions.NotBefore
	if created == "" {
		created = result.Response.Lifetime.Created
	}

	expires := result.Response.Token.Conditions.NotOnOrAfter
	if expires == "" {
		expires = result.Response.Lifetime.Expires
	}

	wresult, err := templates.AdfsSamlTokenTemplate(result.Response.Token.Inner, created, expires, c.RelyingParty)
	if err != nil {
		return "", "", err
	}

	parsedURL, err := url.Parse(c.SiteURL)
	if err != nil {
		return "", "", err
	}

	rootSiteURL := fmt.Sprintf("%s://%s", parsedURL.Scheme, parsedURL.Host)

	params := url.Values{}
	params.Set("wa", "wsignin1.0")
	params.Set("wctx", rootSiteURL+"/_layouts/Authenticate.aspx?Source=%2F")
	params.Set("wresult", wresult)

	// proxyURL, _ := url.Parse("http://127.0.0.1:8888")
	// http.DefaultTransport = &http.Transport{Proxy: http.ProxyURL(proxyURL), TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}

	client = &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	req, err = http.NewRequest("POST", rootSiteURL+"/_trust/", strings.NewReader(params.Encode()))
	if err != nil {
		return "", "", err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if edgeCookie != "" {
		req.Header.Set("Cookie", edgeCookie)
	}

	resp, err = client.Do(req)
	if err != nil {
		return "", "", err
	}
	defer func() { _ = resp.Body.Close() }()

	authCookie := resp.Header.Get("Set-Cookie") // FedAuth
	authCookie = strings.Split(authCookie, ";")[0]

	return authCookie, expires, nil
}

// WAP auth flow - TODO: refactor
func wapAuthFlow(c *AuthCnfg) (string, string, error) {
	client := &http.Client{
		// Disabling redirect so response 302 location can be resolved
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	resp, err := client.Get(c.SiteURL)
	if err != nil {
		return "", "", err
	}
	defer func() { _ = resp.Body.Close() }()

	// Response location with WAP login endpoint is used to send form auth request
	redirect, err := resp.Location()
	if err != nil {
		return "", "", err
	}

	redirectURL := fmt.Sprintf("%s", redirect)

	params := url.Values{}
	params.Set("UserName", c.Username)
	params.Set("Password", c.Password)
	params.Set("AuthMethod", "FormsAuthentication")

	resp, err = client.Post(redirectURL, "application/x-www-form-urlencoded", strings.NewReader(params.Encode()))
	if err != nil {
		return "", "", err
	}
	defer func() { _ = resp.Body.Close() }()

	// Request to redirect URL using MSISAuth
	req, err := http.NewRequest("GET", redirectURL, nil)
	if err != nil {
		return "", "", err
	}
	msisAuthCookie := resp.Header.Get("Set-Cookie")

	if msisAuthCookie == "" {
		err = errors.New("msisAuthCookie is empty, that might be the result of incorrect username and password")
		return "", "", err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/71.0.3578.98 Safari/537.36")
	req.Header.Set("Cookie", msisAuthCookie)

	resp, err = client.Do(req)
	if err != nil {
		return "", "", err
	}
	defer func() { _ = resp.Body.Close() }()

	// Yet another redirect using JWT at this point (spUrl?authToken=JWT&client-request-id=)
	redirect, err = resp.Location()
	if err != nil {
		return "", "", err
	}
	redirectURL = fmt.Sprintf("%s", redirect)

	req, err = http.NewRequest("GET", redirectURL, nil)
	if err != nil {
		return "", "", err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/71.0.3578.98 Safari/537.36")
	// req.Header.Set("Cookie", msisAuthCookie) // brakes it all

	resp, err = client.Do(req)
	if err != nil {
		return "", "", err
	}
	defer func() { _ = resp.Body.Close() }()

	// TODO: get expiry somehow
	authCookie := resp.Header.Get("Set-Cookie") // EdgeAccessCookie
	authCookie = strings.Split(authCookie, ";")[0]

	var fedAuthExpire string

	// ADFS behind WAP scenario, similar to the ordinary ADFS but requires EdgeAccessCookie
	if redirect, err := resp.Location(); err == nil {
		if strings.Contains(fmt.Sprintf("%s", redirect), "/_layouts/15/Authenticate.aspx") {
			redirectURL = fmt.Sprintf("%s", redirect)
			client := &http.Client{}

			req, err = http.NewRequest("GET", redirectURL, nil)
			if err != nil {
				return "", "", err
			}
			req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/71.0.3578.98 Safari/537.36")
			req.Header.Set("Cookie", authCookie)

			resp, err = client.Do(req)
			if err != nil {
				return "", "", err
			}
			defer func() { _ = resp.Body.Close() }()

			cc := *c
			cc.RelyingParty = resp.Request.URL.Query().Get("wtrealm")
			cc.AdfsCookie = "FedAuth"

			fedAuthCookie, expire, err := adfsAuthFlow(&cc, authCookie)
			if err != nil {
				return "", "", err
			}

			fedAuthExpire = expire
			authCookie += "; " + fedAuthCookie
		}
	}

	// fmt.Printf(authCookie)

	return authCookie, fedAuthExpire, nil
}
