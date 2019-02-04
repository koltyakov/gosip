package adfs

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/koltyakov/gosip/cpass"
)

// AuthCnfg : auth config structure
type AuthCnfg struct {
	SiteURL      string `json:"siteUrl"`
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
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &c)

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
func (c *AuthCnfg) SetAuth(req *http.Request) error {
	authCookie, err := c.GetAuth()
	if err != nil {
		return err
	}
	req.Header.Set("Cookie", authCookie)
	return nil
}
