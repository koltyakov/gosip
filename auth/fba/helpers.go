package fba

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/patrickmn/go-cache"

	"github.com/koltyakov/gosip/templates"
)

var (
	storage = cache.New(5*time.Minute, 10*time.Minute)
)

// GetAuth gets authentication
func GetAuth(c *AuthCnfg) (string, int64, error) {
	if c.client == nil {
		c.client = &http.Client{}
	}

	parsedURL, err := url.Parse(c.SiteURL)
	if err != nil {
		return "", 0, err
	}

	cacheKey := parsedURL.Host + "@" + c.GetStrategy() + "@" + c.Username + "@" + c.Password
	if authCookie, exp, found := storage.GetWithExpiration(cacheKey); found {
		return authCookie.(string), exp.Unix(), nil
	}

	endpoint := fmt.Sprintf("%s://%s/_vti_bin/authentication.asmx", parsedURL.Scheme, parsedURL.Host)
	soapBody, err := templates.FbaWsTemplate(c.Username, c.Password)
	if err != nil {
		return "", 0, err
	}

	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer([]byte(soapBody)))
	if err != nil {
		return "", 0, err
	}

	req.Header.Set("Content-Type", "text/xml;charset=utf-8")

	// client := &http.Client{}
	resp, err := c.client.Do(req)
	if err != nil {
		return "", 0, err
	}
	defer func() {
		if resp != nil && resp.Body != nil {
			_ = resp.Body.Close()
		}
	}()

	res, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", 0, err
	}

	// fmt.Printf("FBA: %s\n", string(res))

	type fbaResponse struct {
		ErrorCode      string        `xml:"Body>LoginResponse>LoginResult>ErrorCode"`
		CookieName     string        `xml:"Body>LoginResponse>LoginResult>CookieName"`
		TimeoutSeconds time.Duration `xml:"Body>LoginResponse>LoginResult>TimeoutSeconds"`
	}
	result := &fbaResponse{}
	if err := xml.Unmarshal(res, &result); err != nil {
		return "", 0, err
	}

	if result.ErrorCode != "NoError" {
		return "", 0, errors.New(result.ErrorCode)
	}

	if result.ErrorCode == "PasswordNotMatch" {
		return "", 0, errors.New("password doesn't not match")
	}

	// fmt.Printf("FBA: %s\n", string(result.CookieName))

	authCookie := resp.Header.Get("Set-Cookie") // TODO: parse FBA cookie only (?)
	expiry := (result.TimeoutSeconds - 60) * time.Second
	exp := time.Now().Add(expiry).Unix()

	storage.Set(cacheKey, authCookie, expiry)

	return authCookie, exp, nil
}
