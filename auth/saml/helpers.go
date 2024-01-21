package saml

import (
	"bytes"
	"context"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/patrickmn/go-cache"

	"github.com/koltyakov/gosip/templates"
)

var (
	storage        = cache.New(5*time.Minute, 10*time.Minute)
	loginEndpoints = map[spoEnv]string{
		spoProd:   "login.microsoftonline.com",
		spoGerman: "login.microsoftonline.de",
		spoChina:  "login.chinacloudapi.cn",
		spoUSGov:  "login-us.microsoftonline.com",
		spoUSDef:  "login-us.microsoftonline.com",
	}
)

// GetAuth gets authentication
func GetAuth(ctx context.Context, c *AuthCnfg) (string, int64, error) {
	if c.client == nil {
		c.client = &http.Client{}
	}

	parsedURL, err := url.Parse(c.SiteURL)
	if err != nil {
		return "", 0, err
	}

	cacheKey := parsedURL.Host + "@" + c.GetStrategy() + "@" + c.Username + "@" + c.Password
	if authToken, exp, found := storage.GetWithExpiration(cacheKey); found {
		return authToken.(string), exp.Unix(), nil
	}

	authCookie, notAfter, err := getSecurityToken(ctx, c)
	if err != nil {
		return "", 0, err
	}

	notAfterTime, _ := time.Parse(time.RFC3339, notAfter)
	expiry := time.Until(notAfterTime) - 60*time.Second
	exp := time.Now().Add(expiry).Unix()

	storage.Set(cacheKey, authCookie, expiry)

	return authCookie, exp, nil
}

func getSecurityToken(ctx context.Context, c *AuthCnfg) (string, string, error) {
	if c.client == nil {
		c.client = &http.Client{}
	}

	loginEndpoint := loginEndpoints[resolveSPOEnv(c.SiteURL)]
	endpoint := fmt.Sprintf("https://%s/GetUserRealm.srf", loginEndpoint)

	params := url.Values{}
	params.Set("login", c.Username)

	// client := &http.Client{
	// 	CheckRedirect: func(req *http.Request, via []*http.Request) error {
	// 		return http.ErrUseLastResponse
	// 	},
	// }
	c.client.CheckRedirect = doNotCheckRedirect

	req, err := http.NewRequestWithContext(ctx, "POST", endpoint, strings.NewReader(params.Encode()))
	if err != nil {
		return "", "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := c.client.Do(req)
	if err != nil {
		return "", "", err
	}
	defer func() {
		if resp != nil && resp.Body != nil {
			_ = resp.Body.Close()
		}
	}()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", err
	}

	type userReadlmResponse struct {
		NameSpaceType       string `json:"NameSpaceType"`
		DomainName          string `json:"DomainName"`
		FederationBrandName string `json:"FederationBrandName"`
		CloudInstanceName   string `json:"CloudInstanceName"`
		State               int    `json:"State"`
		UserState           int    `json:"UserState"`
		Login               string `json:"Login"`
		AuthURL             string `json:"AuthURL"`
	}

	userRealm := &userReadlmResponse{}
	if err := json.Unmarshal(data, &userRealm); err != nil {
		return "", "", err
	}

	// fmt.Printf("Results: %v\n", userRealm.NameSpaceType)

	if userRealm.NameSpaceType == "" {
		return "", "", errors.New("unable to define namespace type for Online authentiation")
	}

	if userRealm.NameSpaceType == "Managed" {
		return getSecurityTokenWithOnline(ctx, c)
	}

	if userRealm.NameSpaceType == "Federated" {
		return getSecurityTokenWithAdfs(ctx, userRealm.AuthURL, c)
	}

	return "", "", fmt.Errorf("unable to resolve namespace authentiation type. Type received: %s", userRealm.NameSpaceType)
}

func getSecurityTokenWithOnline(ctx context.Context, c *AuthCnfg) (string, string, error) {
	if c.client == nil {
		c.client = &http.Client{}
	}

	parsedURL, err := url.Parse(c.SiteURL)
	if err != nil {
		return "", "", err
	}

	formsEndpoint := fmt.Sprintf("%s://%s/_forms/default.aspx?wa=wsignin1.0", parsedURL.Scheme, parsedURL.Host)
	samlBody, err := templates.OnlineSamlWsfedTemplate(formsEndpoint, c.Username, c.Password)
	if err != nil {
		return "", "", err
	}

	loginEndpoint := loginEndpoints[resolveSPOEnv(c.SiteURL)]
	stsEndpoint := fmt.Sprintf("https://%s/extSTS.srf", loginEndpoint)

	req, err := http.NewRequestWithContext(ctx, "POST", stsEndpoint, bytes.NewBuffer([]byte(samlBody)))
	if err != nil {
		return "", "", err
	}

	req.Header.Set("Content-Type", "application/soap+xml;charset=utf-8")

	// client := &http.Client{
	// 	CheckRedirect: func(req *http.Request, via []*http.Request) error {
	// 		return http.ErrUseLastResponse
	// 	},
	// }
	c.client.CheckRedirect = doNotCheckRedirect

	resp, err := c.client.Do(req)
	if err != nil {
		return "", "", err
	}
	defer func() {
		if resp != nil && resp.Body != nil {
			_ = resp.Body.Close()
		}
	}()

	xmlResponse, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", err
	}

	type samlAssertion struct {
		Fault    string `xml:"Body>Fault>Reason>Text"`
		Response struct {
			BinaryToken string `xml:"RequestedSecurityToken>BinarySecurityToken"`
			Lifetime    struct {
				Created string `xml:"Created"`
				Expires string `xml:"Expires"`
			} `xml:"Lifetime"`
		} `xml:"Body>RequestSecurityTokenResponse"`
	}
	result := &samlAssertion{}
	if err := xml.Unmarshal(xmlResponse, &result); err != nil {
		return "", "", err
	}

	// fmt.Printf("BinaryToken, %s\n", result.Response.BinaryToken)

	resp, err = c.client.Post(formsEndpoint, "application/x-www-form-urlencoded", strings.NewReader(result.Response.BinaryToken))
	if err != nil {
		return "", "", err
	}
	defer func() {
		if resp != nil && resp.Body != nil {
			_ = resp.Body.Close()
		}
	}()

	if _, err := io.Copy(io.Discard, resp.Body); err != nil {
		return "", "", err
	}

	// cookie := resp.Header.Get("Set-Cookie") // TODO: parse FedAuth and rtFa cookies only (?)
	// fmt.Printf("Cookie: %s\n", cookie)
	// fmt.Printf("Resp2, %v\n", resp.StatusCode)

	var authCookie string
	for _, coo := range resp.Cookies() {
		if coo.Name == "rtFa" || coo.Name == "FedAuth" {
			authCookie += coo.String() + "; "
		}
	}

	return authCookie, result.Response.Lifetime.Expires, nil
}

// TODO: test the method, it possibly contains issues and extra complexity
func getSecurityTokenWithAdfs(ctx context.Context, adfsURL string, c *AuthCnfg) (string, string, error) {
	if c.client == nil {
		c.client = &http.Client{}
	}

	parsedAdfsURL, err := url.Parse(adfsURL)
	if err != nil {
		return "", "", err
	}

	// proxyURL, _ := url.Parse("http://127.0.0.1:8888")
	// http.DefaultTransport = &http.Transport{Proxy: http.ProxyURL(proxyURL), TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}

	usernameMixedURL := fmt.Sprintf("%s://%s/adfs/services/trust/13/usernamemixed", parsedAdfsURL.Scheme, parsedAdfsURL.Host)
	samlBody, err := templates.AdfsSamlWsfedTemplate(usernameMixedURL, c.Username, c.Password, "urn:federation:MicrosoftOnline")
	if err != nil {
		return "", "", err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", usernameMixedURL, bytes.NewBuffer([]byte(samlBody)))
	if err != nil {
		return "", "", err
	}

	req.Header.Set("Content-Type", "application/soap+xml;charset=utf-8")

	// client := &http.Client{}
	c.client.CheckRedirect = doNotCheckRedirect
	resp, err := c.client.Do(req)
	if err != nil {
		return "", "", err
	}
	defer func() {
		if resp != nil && resp.Body != nil {
			_ = resp.Body.Close()
		}
	}()

	res, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", err
	}

	// fmt.Printf("ADFS: %s\n", string(res))

	type samlAssertion struct {
		Response struct {
			Fault string `xml:"Fault>Reason>Text"`
			Token struct {
				Inner      []byte `xml:",innerxml"`
				Conditions struct {
					NotBefore    string `xml:"NotBefore,attr"`
					NotOnOrAfter string `xml:"NotOnOrAfter,attr"`
				} `xml:"Assertion>Conditions"`
			} `xml:"RequestSecurityTokenResponseCollection>RequestSecurityTokenResponse>RequestedSecurityToken"`
		} `xml:"Body"`
	}

	result := &samlAssertion{}
	if err := xml.Unmarshal(res, &result); err != nil {
		return "", "", err
	}

	// fmt.Printf("Token: %s", result.Response.Token.Inner)

	if result.Response.Fault != "" {
		return "", "", errors.New(result.Response.Fault)
	}

	// parsedURL, err := url.Parse(adfsURL)
	parsedURL, err := url.Parse(c.SiteURL)
	if err != nil {
		return "", "", err
	}

	rootSite := fmt.Sprintf("%s://%s", parsedURL.Scheme, parsedURL.Host)
	tokenRequest, err := templates.OnlineSamlWsfedAdfsTemplate(rootSite, string(result.Response.Token.Inner))
	if err != nil {
		return "", "", err
	}

	// fmt.Printf("tokenRequest: %s\n", tokenRequest)

	stsEndpoint := "https://login.microsoftonline.com/extSTS.srf" // TODO: mapping

	req, err = http.NewRequestWithContext(ctx, "POST", stsEndpoint, bytes.NewBuffer([]byte(tokenRequest)))
	if err != nil {
		return "", "", err
	}

	req.Header.Set("Content-Type", "application/soap+xml;charset=utf-8")

	resp, err = c.client.Do(req)
	if err != nil {
		return "", "", err
	}
	defer func() {
		if resp != nil && resp.Body != nil {
			_ = resp.Body.Close()
		}
	}()

	xmlResponse, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", err
	}

	// fmt.Printf("token: %s\n", xmlResponse)

	type tokenAssertion struct {
		Fault    string `xml:"Body>Fault>Reason>Text"`
		Response struct {
			BinaryToken string `xml:"RequestedSecurityToken>BinarySecurityToken"`
			Lifetime    struct {
				Created string `xml:"Created"`
				Expires string `xml:"Expires"`
			} `xml:"Lifetime"`
		} `xml:"Body>RequestSecurityTokenResponse"`
	}

	tokenResult := &tokenAssertion{}
	if err := xml.Unmarshal(xmlResponse, &tokenResult); err != nil {
		return "", "", err
	}

	if tokenResult.Response.BinaryToken == "" {
		return "", "", errors.New("can't extract binary token")
	}

	// client = &http.Client{
	// 	CheckRedirect: func(req *http.Request, via []*http.Request) error {
	// 		return http.ErrUseLastResponse
	// 	},
	// }
	c.client.CheckRedirect = doNotCheckRedirect

	formsEndpoint := fmt.Sprintf("%s://%s/_forms/default.aspx?wa=wsignin1.0", parsedURL.Scheme, parsedURL.Host)
	resp, err = c.client.Post(formsEndpoint, "application/x-www-form-urlencoded", strings.NewReader(tokenResult.Response.BinaryToken))
	if err != nil {
		return "", "", err
	}
	defer func() {
		if resp != nil && resp.Body != nil {
			_ = resp.Body.Close()
		}
	}()

	if _, err := io.Copy(io.Discard, resp.Body); err != nil {
		return "", "", err
	}

	var authCookie string
	for _, coo := range resp.Cookies() {
		if coo.Name == "rtFa" || coo.Name == "FedAuth" {
			authCookie += coo.String() + "; "
		}
	}

	return authCookie, tokenResult.Response.Lifetime.Expires, nil
}

// doNotCheckRedirect *http.Client CheckRedirect callback to ignore redirects
func doNotCheckRedirect(_ *http.Request, _ []*http.Request) error {
	return http.ErrUseLastResponse
}
