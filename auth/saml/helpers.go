package saml

import (
	"bytes"
	"encoding/json"
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

	cacheKey := parsedURL.Host + "@saml@" + creds.Username + "@" + creds.Password
	if authToken, found := storage.Get(cacheKey); found {
		return authToken.(string), nil
	}

	authCookie, notAfter, err := getSecurityToken(creds)
	if err != nil {
		return "", nil
	}

	notAfterTime, _ := time.Parse(time.RFC3339, notAfter)
	storage.Set(cacheKey, authCookie, (time.Until(notAfterTime)-60)*time.Second)

	return authCookie, nil
}

func getSecurityToken(creds *AuthCnfg) (string, string, error) {
	endpoint := "https://login.microsoftonline.com/GetUserRealm.srf" // TODO: endpoints mapping

	params := url.Values{}
	params.Set("login", creds.Username)

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	resp, err := client.Post(endpoint, "application/x-www-form-urlencoded", strings.NewReader(params.Encode()))
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
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
	}

	userRealm := &userReadlmResponse{}
	if err := json.Unmarshal(data, &userRealm); err != nil {
		return "", "", err
	}

	// fmt.Printf("Results: %v\n", userRealm.NameSpaceType)

	if userRealm.NameSpaceType == "" {
		return "", "", errors.New("Unable to define namespace type for Online authentiation")
	}

	if userRealm.NameSpaceType == "Managed" {
		return getSecurityTokenWithOnline(creds)
	}

	if userRealm.NameSpaceType == "Federated" {
		// return getSecurityTokenWithAdfs(userRealm.AuthURL)
	}

	return "", "", fmt.Errorf("Unable to resolve namespace authentiation type. Type received: %s", userRealm.NameSpaceType)
}

func getSecurityTokenWithOnline(creds *AuthCnfg) (string, string, error) {
	parsedURL, err := url.Parse(creds.SiteURL)
	if err != nil {
		return "", "", err
	}

	formsEndpoint := fmt.Sprintf("%s://%s/_forms/default.aspx?wa=wsignin1.0", parsedURL.Scheme, parsedURL.Host)
	samlBody, err := templates.OnlineSamlWsfedTemplate(formsEndpoint, creds.Username, creds.Password)
	if err != nil {
		return "", "", err
	}

	stsEndpoint := "https://login.microsoftonline.com/extSTS.srf" // TODO: add mapping for diff SPOs

	req, err := http.NewRequest("POST", stsEndpoint, bytes.NewBuffer([]byte(samlBody)))
	if err != nil {
		return "", "", err
	}

	req.Header.Set("Content-Type", "application/soap+xml;charset=utf-8")

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	xmlResponse, err := ioutil.ReadAll(resp.Body)
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

	resp, err = client.Post(formsEndpoint, "application/x-www-form-urlencoded", strings.NewReader(result.Response.BinaryToken))
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

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

// func getSecurityTokenWithAdfs() {
// 	// Reuse ADFS helpers
// }
