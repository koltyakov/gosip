// Package azurecert implements AAD Certificate Auth Flow
// See more:
//   - https://docs.microsoft.com/en-us/azure/developer/go/azure-sdk-authorization#use-file-based-authentication
//
// Amongst supported platform versions are:
//   - SharePoint Online + Azure
package azurecert

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
	"time"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure/auth"
	"github.com/koltyakov/gosip"
	"github.com/koltyakov/gosip/cpass"
	"github.com/patrickmn/go-cache"
)

var (
	storage = cache.New(5*time.Minute, 10*time.Minute)
)

// AuthCnfg - AAD Certificate Auth Flow
/* Config sample:
{
	"siteUrl": "https://contoso.sharepoint.com/sites/test",
	"tenantId": "e4d43069-8ecb-49c4-8178-5bec83c53e9d",
	"clientId": "628cc712-c9a4-48f0-a059-af64bdbb4be5",
	"certPath": "cert.pfx",
	"certPass": "password"
}
*/
type AuthCnfg struct {
	SiteURL  string `json:"siteUrl"`  // SPSite or SPWeb URL, which is the context target for the API calls
	TenantID string `json:"tenantId"` // Azure Tenant ID
	ClientID string `json:"clientId"` // Azure Client ID
	CertPath string `json:"certPath"` // Azure certificate (.pfx) file location, relative to config location or absolute
	CertPass string `json:"certPass"` // Azure certificate export password

	authorizer  autorest.Authorizer
	privateFile string
	masterKey   string
}

// ReadConfig reads private config with auth options
func (c *AuthCnfg) ReadConfig(privateFile string) error {
	c.privateFile = privateFile
	f, err := os.Open(privateFile)
	if err != nil {
		return err
	}
	defer func() { _ = f.Close() }()
	byteValue, _ := io.ReadAll(f)
	return c.ParseConfig(byteValue)
}

// ParseConfig parses credentials from a provided JSON byte array content
func (c *AuthCnfg) ParseConfig(byteValue []byte) error {
	if err := json.Unmarshal(byteValue, &c); err != nil {
		return err
	}
	c.CertPath = path.Join(path.Dir(c.privateFile), c.CertPath)
	crypt := cpass.Cpass(c.masterKey)
	secret, err := crypt.Decode(c.CertPass)
	if err == nil {
		c.CertPass = secret
	}
	return nil
}

// WriteConfig writes private config with auth options
func (c *AuthCnfg) WriteConfig(privateFile string) error {
	crypt := cpass.Cpass(c.masterKey)
	secret, err := crypt.Encode(c.CertPass)
	if err != nil {
		return err
	}
	config := &AuthCnfg{
		SiteURL:  c.SiteURL,
		TenantID: c.TenantID,
		ClientID: c.ClientID,
		CertPath: c.CertPass,
		CertPass: secret,
	}
	file, _ := json.MarshalIndent(config, "", "  ")
	return os.WriteFile(privateFile, file, 0644)
}

// SetMasterkey defines custom masterkey
func (c *AuthCnfg) SetMasterkey(masterKey string) { c.masterKey = masterKey }

// GetAuth authenticates, receives access token
func (c *AuthCnfg) GetAuth() (string, int64, error) {
	if c.authorizer == nil {
		u, _ := url.Parse(c.SiteURL)
		resource := fmt.Sprintf("https://%s", u.Host)

		config := auth.NewClientCertificateConfig(c.CertPath, c.CertPass, c.ClientID, c.TenantID)
		config.Resource = resource

		authorizer, err := config.Authorizer()
		if err != nil {
			return "", 0, err
		}
		c.authorizer = authorizer
	}

	// token, err := config.ServicePrincipalToken()
	// if err != nil {
	// 	return "", 0, err
	// }
	// return token.Token().AccessToken, token.Token().Expires().Unix(), nil

	return c.getToken()
}

// GetSiteURL gets SharePoint siteURL
func (c *AuthCnfg) GetSiteURL() string { return c.SiteURL }

// GetStrategy gets auth strategy name
func (c *AuthCnfg) GetStrategy() string { return "azurecert" }

// SetAuth authenticates request
// noinspection GoUnusedParameter
func (c *AuthCnfg) SetAuth(req *http.Request, httpClient *gosip.SPClient) error {
	authToken, _, err := c.GetAuth()
	if err != nil {
		return err
	}
	// _, err := c.authorizer.WithAuthorization()(preparer{}).Prepare(req)
	req.Header.Set("Authorization", "Bearer "+authToken)
	return err
}

// Getting token with prepare for external usage scenarious
func (c *AuthCnfg) getToken() (string, int64, error) {
	// Get from cache
	parsedURL, err := url.Parse(c.SiteURL)
	if err != nil {
		return "", 0, err
	}
	cacheKey := parsedURL.Host + "@" + c.GetStrategy() + "@" + c.TenantID + "@" + c.ClientID
	if accessToken, exp, found := storage.GetWithExpiration(cacheKey); found {
		return accessToken.(string), exp.Unix(), nil
	}

	// Get token
	req, _ := http.NewRequest("GET", c.SiteURL, nil)
	req, err = c.authorizer.WithAuthorization()(preparer{}).Prepare(req)
	if err != nil {
		return "", 0, err
	}
	token := strings.Replace(req.Header.Get("Authorization"), "Bearer ", "", 1)
	tt := strings.Split(token, ".")
	if len(tt) != 3 {
		return "", 0, fmt.Errorf("incorrect jwt")
	}
	jsonBytes, err := base64.RawURLEncoding.DecodeString(tt[1])
	if err != nil {
		return "", 0, fmt.Errorf("can't decode jwt base64 string")
	}
	j := struct {
		Exp int64 `json:"exp"`
	}{}
	_ = json.Unmarshal(jsonBytes, &j)

	// Save to cache
	exp := time.Unix(j.Exp, 0).Add(-60 * time.Second)
	storage.Set(cacheKey, token, time.Until(exp))

	// fmt.Println(time.Until(exp))

	return token, exp.Unix(), nil
}

// Preparer implements autorest.Preparer interface
type preparer struct{}

// Prepare satisfies autorest.Preparer interface
func (p preparer) Prepare(req *http.Request) (*http.Request, error) { return req, nil }
