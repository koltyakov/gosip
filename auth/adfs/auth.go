/*
Package adfs implements ADFS Auth (user credentials authentication)

Amongst supported platform versions are:
	- SharePoint Online (SPO)
	- On-Prem: 2019, 2016, and 2013
*/
package adfs

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/koltyakov/gosip"
	"github.com/koltyakov/gosip/cpass"
)

// AuthCnfg - ADFS auth config structure
/* On-Premises config sample:
{
  "siteUrl": "https://www.contoso.com/sites/test",
  "username": "john.doe@contoso.com",
  "password": "this-is-not-a-real-password",
  "relyingParty": "urn:sharepoint:www",
  "adfsUrl": "https://login.contoso.com",
  "adfsCookie": "FedAuth"
}
*/
/* On-Premises behind WAP config sample:
{
  "siteUrl": "https://www.contoso.com/sites/test",
  "username": "john.doe@contoso.com",
  "password": "this-is-not-a-real-password",
  "relyingParty": "urn:AppProxy:com",
  "adfsUrl": "https://login.contoso.com",
  "adfsCookie": "EdgeAccessCookie"
}
*/
/* SharePoint Online config sample:
{
  "siteUrl": "https://www.contoso.com/sites/test",
  "username": "john.doe@contoso.com",
  "password": "this-is-not-a-real-password"
}
*/
type AuthCnfg struct {
	SiteURL      string `json:"siteUrl"` // SPSite or SPWeb URL, which is the context target for the API calls
	Domain       string `json:"domain"`
	Username     string `json:"username"`
	Password     string `json:"password"`
	RelyingParty string `json:"relyingParty"`
	AdfsURL      string `json:"adfsUrl"`
	AdfsCookie   string `json:"adfsCookie"`

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

	if c.Domain != "" && !strings.Contains(c.Username, "\\") && !strings.Contains(c.Username, "@") {
		c.Username = c.Domain + "\\" + c.Username
	}

	if c.AdfsCookie == "" {
		c.AdfsCookie = "FedAuth"
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
		SiteURL:      c.SiteURL,
		Username:     c.Username,
		Domain:       c.Domain,
		Password:     pass,
		RelyingParty: c.RelyingParty,
		AdfsURL:      c.AdfsURL,
		AdfsCookie:   c.AdfsCookie,
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
	return "adfs"
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
