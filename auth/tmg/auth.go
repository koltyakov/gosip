/*
Package tmg implements FBA authentication behind TMG (Microsoft Forefront Threat Management Gateway)

Currently is legacy but was a popular way of exposing SharePoint into external world back in the days.

Amongst supported platform versions are:
	- On-Prem: 2019, 2016, and 2013
*/
package tmg

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/koltyakov/gosip"
	"github.com/koltyakov/gosip/cpass"
)

// AuthCnfg - FBA bihind TMG auth config structure
/* On-Premises config sample:
{
  "siteUrl": "https://www.contoso.com/sites/test",
  "username": "john.doe@contoso.com",
  "password": "this-is-not-a-real-password"
}
*/
type AuthCnfg struct {
	SiteURL  string `json:"siteUrl"` // SPSite or SPWeb URL, which is the context target for the API calls
	Username string `json:"username"`
	Password string `json:"password"`

	masterKey string
}

// ReadConfig : reads private config with auth options
func (c *AuthCnfg) ReadConfig(privateFile string) error {
	jsonFile, err := os.Open(privateFile)
	if err != nil {
		return err
	}
	defer func() { _ = jsonFile.Close() }()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	if err := json.Unmarshal(byteValue, &c); err != nil {
		return err
	}

	crypt := cpass.Cpass(c.masterKey)
	pass, err := crypt.Decode(c.Password)
	if err == nil {
		c.Password = pass
	}

	return nil
}

// WriteConfig : writes private config with auth options
func (c *AuthCnfg) WriteConfig(privateFile string) error {
	crypt := cpass.Cpass(c.masterKey)
	pass, err := crypt.Encode(c.Password)
	if err != nil {
		pass = c.Password
	}
	config := &AuthCnfg{
		SiteURL:  c.SiteURL,
		Username: c.Username,
		Password: pass,
	}
	file, _ := json.MarshalIndent(config, "", "  ")
	return ioutil.WriteFile(privateFile, file, 0644)
}

// SetMasterkey : defines custom masterkey
func (c *AuthCnfg) SetMasterkey(masterKey string) {
	c.masterKey = masterKey
}

// GetAuth : authenticates, receives access token
func (c *AuthCnfg) GetAuth() (string, error) {
	return GetAuth(c)
}

// GetSiteURL : gets siteURL
func (c *AuthCnfg) GetSiteURL() string {
	return c.SiteURL
}

// GetStrategy : gets auth strategy name
func (c *AuthCnfg) GetStrategy() string {
	return "tmg"
}

// SetAuth : authenticate request
//noinspection GoUnusedParameter
func (c *AuthCnfg) SetAuth(req *http.Request, httpClient *gosip.SPClient) error {
	authCookie, err := c.GetAuth()
	if err != nil {
		return err
	}
	req.Header.Set("Cookie", authCookie)
	return nil
}
