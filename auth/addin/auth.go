package addin

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/koltyakov/gosip"
	"github.com/koltyakov/gosip/cpass"
)

// AuthCnfg : auth config structure
type AuthCnfg struct {
	SiteURL      string `json:"siteUrl"`
	ClientID     string `json:"clientId"`
	ClientSecret string `json:"clientSecret"`
	Realm        string `json:"realm"`

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
	json.Unmarshal(byteValue, &c)

	crypt := cpass.Cpass(c.masterKey)
	secret, err := crypt.Decode(c.ClientSecret)
	if err == nil {
		c.ClientSecret = secret
	}

	return nil
}

// WriteConfig : writes private config with auth options
func (c *AuthCnfg) WriteConfig(privateFile string) error {
	crypt := cpass.Cpass(c.masterKey)
	secret, err := crypt.Encode(c.ClientSecret)
	if err != nil {
		secret = c.ClientSecret
	}
	conf := &AuthCnfg{
		SiteURL:      c.SiteURL,
		ClientID:     c.ClientID,
		ClientSecret: secret,
		Realm:        c.Realm,
	}
	file, _ := json.MarshalIndent(conf, "", "  ")
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
	return "addin"
}

// SetAuth : authenticate request
func (c *AuthCnfg) SetAuth(req *http.Request, httpClient *gosip.SPClient) error {
	authToken, err := c.GetAuth()
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+authToken)
	return nil
}
