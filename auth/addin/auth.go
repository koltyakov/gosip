/*
Package addin implements AddIn Only Auth

This type of authentication uses AddIn Only policy and OAuth bearer tokens for authenticating HTTP requests.

Amongst supported platform versions are:
	- SharePoint Online (SPO)
*/
package addin

import (
	"encoding/json"
	"io"
	"net/http"
	"os"

	"github.com/koltyakov/gosip"
	"github.com/koltyakov/gosip/cpass"
)

// AuthCnfg - AddIn Only auth config structure
/* SharePoint Online config sample:
{
  "siteUrl": "https://contoso.sharepoint.com/sites/test",
  "clientId": "e2763c6d-7ee6-41d6-b15c-dd1f75f90b8f",
  "clientSecret": "OqDSAAuBChzI+uOX0OUhXxiOYo1g6X7mjXCVA9mSF/0="
}
*/
type AuthCnfg struct {
	SiteURL      string `json:"siteUrl"`      // SPSite or SPWeb URL, which is the context target for the API calls
	ClientID     string `json:"clientId"`     // Client ID obtained when registering the AddIn
	ClientSecret string `json:"clientSecret"` // Client Secret obtained when registering the AddIn
	Realm        string `json:"realm"`        // Your SharePoint Online tenant ID (optional)

	masterKey string
	client    *http.Client
}

// ReadConfig reads private config with auth options
func (c *AuthCnfg) ReadConfig(privateFile string) error {
	jsonFile, err := os.Open(privateFile)
	if err != nil {
		return err
	}
	defer func() { _ = jsonFile.Close() }()

	byteValue, _ := io.ReadAll(jsonFile)
	return c.ParseConfig(byteValue)
}

// ParseConfig parses credentials from a provided JSON byte array content
func (c *AuthCnfg) ParseConfig(byteValue []byte) error {
	if err := json.Unmarshal(byteValue, &c); err != nil {
		return err
	}

	crypt := cpass.Cpass(c.masterKey)
	secret, err := crypt.Decode(c.ClientSecret)
	if err == nil {
		c.ClientSecret = secret
	}

	return nil
}

// WriteConfig writes private config with auth options
func (c *AuthCnfg) WriteConfig(privateFile string) error {
	crypt := cpass.Cpass(c.masterKey)
	secret, err := crypt.Encode(c.ClientSecret)
	if err != nil {
		secret = c.ClientSecret
	}
	config := &AuthCnfg{
		SiteURL:      c.SiteURL,
		ClientID:     c.ClientID,
		ClientSecret: secret,
		Realm:        c.Realm,
	}
	file, _ := json.MarshalIndent(config, "", "  ")
	return os.WriteFile(privateFile, file, 0644)
}

// SetMasterkey defines custom masterkey
func (c *AuthCnfg) SetMasterkey(masterKey string) { c.masterKey = masterKey }

// GetAuth authenticates, receives access token
func (c *AuthCnfg) GetAuth() (string, int64, error) { return GetAuth(c) }

// GetSiteURL gets siteURL
func (c *AuthCnfg) GetSiteURL() string { return c.SiteURL }

// GetStrategy gets auth strategy name
func (c *AuthCnfg) GetStrategy() string { return "addin" }

// SetAuth authenticate request
// noinspection GoUnusedParameter
func (c *AuthCnfg) SetAuth(req *http.Request, httpClient *gosip.SPClient) error {
	authToken, _, err := c.GetAuth()
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+authToken)
	return nil
}
