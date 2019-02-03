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

	token, notBefore, notAfter, err := getSamlAssertion(creds)
	if err != nil {
		return "", err
	}

	authCookie, err := postTokenData(token, notBefore, notAfter, creds)
	if err != nil {
		fmt.Printf("Post token error: %v\n", err)
		return "", err
	}

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
	samlBody, err := buildAdfsSamlWsfedTemplate(usernameMixedURL, creds.Username, creds.Password, creds.RelyingParty)
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
		} `xml:"Body>RequestSecurityTokenResponseCollection>RequestSecurityTokenResponse"`
	}
	result := &samlAssertion{}
	if err := xml.Unmarshal(res, &result); err != nil {
		return []byte(""), "", "", err
	}

	if result.Fault != "" {
		return []byte(""), "", "", errors.New(result.Fault)
	}

	return result.Response.Token.Inner, result.Response.Token.Conditions.NotBefore, result.Response.Token.Conditions.NotOnOrAfter, nil
}

func postTokenData(token []byte, notBefore, notAfter string, creds *AuthCnfg) (string, error) {
	wresult, err := buildAdfsSamlTokenTemplate(token, notBefore, notAfter, creds.RelyingParty)
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
