/*
Package ntlm implements NTLM Auth (NTLM handshake)

This type of authentication uses HTTP NTLM handshake in order to obtain authentication header.

Amongst supported platform versions are:
	- On-Prem: 2019, 2016, and 2013
*/
package ntlm

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	ntlmssp "github.com/Azure/go-ntlmssp"
	"github.com/koltyakov/gosip"
	"github.com/koltyakov/gosip/cpass"
)

// AuthCnfg - NTML auth config structure
/* On-Premises config sample:
{
  "siteUrl": "https://www.contoso.com/sites/test",
  "username": "contoso\\john.doe",
  "password": "this-is-not-a-real-password"
}
or
{
  "siteUrl": "https://www.contoso.com/sites/test",
	"username": "john.doe",
	"domain": "contoso",
  "password": "this-is-not-a-real-password"
}
*/
type AuthCnfg struct {
	SiteURL  string `json:"siteUrl"`  // SPSite or SPWeb URL, which is the context target for the API calls
	Domain   string `json:"domain"`   // AD domain name (optional)
	Username string `json:"username"` // AD user name
	Password string `json:"password"` // AD user password

	masterKey string
}

// ReadConfig : reads private config with auth options
func (c *AuthCnfg) ReadConfig(privateFile string) error {
	jsonFile, err := os.Open(privateFile)
	if err != nil {
		return err
	}
	defer jsonFile.Close()

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
		Domain:   c.Domain,
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
	return "", nil
}

// GetSiteURL : gets siteURL
func (c *AuthCnfg) GetSiteURL() string {
	return c.SiteURL
}

// GetStrategy : gets auth strategy name
func (c *AuthCnfg) GetStrategy() string {
	return "ntlm"
}

// SetAuth : authenticate request
func (c *AuthCnfg) SetAuth(req *http.Request, httpClient *gosip.SPClient) error {
	// NTML + Negotiation
	if httpClient.Transport != nil {
		httpClient.Transport = ntlmssp.Negotiator{
			RoundTripper: httpClient.Transport, // custom transport
		}
	} else {
		httpClient.Transport = ntlmssp.Negotiator{
			RoundTripper: &http.Transport{},
		}
	}
	req.SetBasicAuth(c.Username, c.Password)
	return nil
}
