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

	"github.com/koltyakov/gosip/templates"
	cache "github.com/patrickmn/go-cache"
)

var (
	storage = cache.New(5*time.Minute, 10*time.Minute)
)

// GetAuth : get auth
func GetAuth(creds *AuthCnfg) (string, error) {
	parsedURL, err := url.Parse(creds.SiteURL)
	if err != nil {
		return "", err
	}

	cacheKey := parsedURL.Host + "@adfs@" + creds.Username + "@" + creds.Password
	if authCookie, found := storage.Get(cacheKey); found {
		return authCookie.(string), nil
	}

	// In case of WAP
	if creds.AdfsCookie == "EdgeAccessCookie" {
		authCookie, err := wapAuthFlow(creds)
		if err != nil {
			return "", err
		}
		storage.Set(cacheKey, authCookie, time.Duration(30)*time.Minute)
		return authCookie, nil
	}

	token, notBefore, notAfter, err := getSamlAssertion(creds)
	if err != nil {
		return "", err
	}

	authCookie, err := postTokenData(token, notBefore, notAfter, creds)
	if err != nil {
		fmt.Printf("Post token error: %v\n", err)
		return "", err
	}

	// fmt.Printf("Auth Cookie: %s\n", authCookie)

	notAfterTime, _ := time.Parse(time.RFC3339, notAfter)
	storage.Set(cacheKey, authCookie, (time.Until(notAfterTime)-60)*time.Second)

	return authCookie, nil
}

func getSamlAssertion(creds *AuthCnfg) ([]byte, string, string, error) {
	parsedAdfsURL, err := url.Parse(creds.AdfsURL)
	if err != nil {
		return []byte(""), "", "", err
	}

	usernameMixedURL := fmt.Sprintf("%s://%s/adfs/services/trust/13/usernamemixed", parsedAdfsURL.Scheme, parsedAdfsURL.Host)
	samlBody, err := templates.AdfsSamlWsfedTemplate(usernameMixedURL, creds.Username, creds.Password, creds.RelyingParty)
	if err != nil {
		return []byte(""), "", "", err
	}

	req, err := http.NewRequest("POST", usernameMixedURL, bytes.NewBuffer([]byte(samlBody)))
	if err != nil {
		return []byte(""), "", "", err
	}

	req.Header.Set("Content-Type", "application/soap+xml;charset=utf-8")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return []byte(""), "", "", err
	}
	defer resp.Body.Close()

	res, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte(""), "", "", err
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
		return []byte(""), "", "", err
	}

	if result.Fault != "" {
		return []byte(""), "", "", errors.New(result.Fault)
	}

	created := result.Response.Token.Conditions.NotBefore
	if created == "" {
		created = result.Response.Lifetime.Created
	}

	expires := result.Response.Token.Conditions.NotOnOrAfter
	if expires == "" {
		expires = result.Response.Lifetime.Expires
	}

	return result.Response.Token.Inner, created, expires, nil
}

func postTokenData(token []byte, notBefore, notAfter string, creds *AuthCnfg) (string, error) {
	wresult, err := templates.AdfsSamlTokenTemplate(token, notBefore, notAfter, creds.RelyingParty)
	if err != nil {
		return "", err
	}

	parsedURL, err := url.Parse(creds.SiteURL)
	if err != nil {
		return "", err
	}

	rootSiteURL := fmt.Sprintf("%s://%s", parsedURL.Scheme, parsedURL.Host)

	params := url.Values{}
	params.Set("wa", "wsignin1.0")
	params.Set("wctx", rootSiteURL+"/_layouts/Authenticate.aspx?Source=%2F")
	params.Set("wresult", wresult)

	// proxyURL, _ := url.Parse("http://127.0.0.1:8888")
	// http.DefaultTransport = &http.Transport{Proxy: http.ProxyURL(proxyURL), TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	resp, err := client.Post(rootSiteURL+"/_trust/", "application/x-www-form-urlencoded", strings.NewReader(params.Encode()))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	cookie := resp.Header.Get("Set-Cookie") // TODO: parse ADFS cookie only (?)

	return cookie, nil
}

// WAP auth flow - TODO: refactor
func wapAuthFlow(creds *AuthCnfg) (string, error) {
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	resp, err := client.Get(creds.SiteURL)
	if err != nil {
		return "", err
	}

	redirect, err := resp.Location()
	if err != nil {
		return "", err
	}

	redirectURL := fmt.Sprintf("%s", redirect)

	params := url.Values{}
	params.Set("UserName", creds.Username)
	params.Set("Password", creds.Password)
	params.Set("AuthMethod", "FormsAuthentication")

	resp, err = client.Post(redirectURL, "application/x-www-form-urlencoded", strings.NewReader(params.Encode()))
	if err != nil {
		return "", err
	}
	// defer resp.Body.Close()

	tempCookie := resp.Header.Get("Set-Cookie")

	req, err := http.NewRequest("GET", redirectURL, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/71.0.3578.98 Safari/537.36")
	req.Header.Set("Cookie", tempCookie)

	resp, err = client.Do(req)
	if err != nil {
		return "", err
	}

	redirect, err = resp.Location()
	if err != nil {
		return "", err
	}
	redirectURL = fmt.Sprintf("%s", redirect)

	req, err = http.NewRequest("GET", redirectURL, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/71.0.3578.98 Safari/537.36")

	resp, err = client.Do(req)
	if err != nil {
		return "", err
	}

	// TODO: get expirity somehow
	authCookie := resp.Header.Get("Set-Cookie")
	authCookie = strings.Split(authCookie, ";")[0]

	// fmt.Printf(authCookie)

	return authCookie, nil
}
