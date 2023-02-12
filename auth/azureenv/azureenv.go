// Package azureenv implements AAD Environment-Based Auth Flow
// See more:
//   - https://docs.microsoft.com/en-us/azure/developer/go/azure-sdk-authorization#use-environment-based-authentication
//   - https://docs.microsoft.com/en-us/sharepoint/dev/solution-guidance/security-apponly-azuread
//   - https://docs.microsoft.com/en-us/azure/active-directory/develop/howto-create-service-principal-portal
// Amongst supported platform versions are:
//   - SharePoint Online + Azure
// Azure Environment-Based supported strategies:
//   - Client credentials (might not work with SharePoint but require a Certificate-based auth)
//   - Certificate
//   - Username/Password
//   - Managed identity
package azureenv

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure/auth"
	"github.com/koltyakov/gosip"
	"github.com/koltyakov/gosip/cpass"
)

// AuthCnfg - AAD Environment-Based Auth Flow
// No creds settings in auth config are required, the configuration is happening through environment variables:
// https://docs.microsoft.com/en-us/azure/developer/go/azure-sdk-authorization#use-environment-based-authentication
/* Config sample:
{ "siteUrl": "https://contoso.sharepoint.com/sites/test" }
*/
type AuthCnfg struct {
	SiteURL string            `json:"siteUrl"` // SPSite or SPWeb URL, which is the context target for the API calls
	Env     map[string]string `json:"env"`     // AZURE_ environment variables

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
	byteValue, _ := ioutil.ReadAll(f)
	return c.ParseConfig(byteValue)
}

// ParseConfig parses credentials from a provided JSON byte array content
func (c *AuthCnfg) ParseConfig(byteValue []byte) error {
	if err := json.Unmarshal(byteValue, &c); err != nil {
		return err
	}
	crypt := cpass.Cpass(c.masterKey)
	for key, val := range c.Env {
		if key == "AZURE_AUTH_LOCATION" || key == "AZURE_CERTIFICATE_PATH" {
			c.Env[key] = path.Join(path.Dir(c.privateFile), val)
		}
		if strings.Contains(strings.ToLower(key), "_password") || strings.Contains(strings.ToLower(key), "_secret") {
			secret, err := crypt.Decode(val)
			if err == nil {
				c.Env[key] = secret
			}
		}
	}
	return nil
}

// WriteConfig writes private config with auth options
func (c *AuthCnfg) WriteConfig(privateFile string) error {
	config := &AuthCnfg{SiteURL: c.SiteURL}
	file, _ := json.MarshalIndent(config, "", "  ")
	return ioutil.WriteFile(privateFile, file, 0644)
}

// SetMasterkey defines custom masterkey
func (c *AuthCnfg) SetMasterkey(masterKey string) { c.masterKey = masterKey }

// GetAuth authenticates, receives access token
func (c *AuthCnfg) GetAuth() (string, int64, error) {
	if c.authorizer == nil {
		u, _ := url.Parse(c.SiteURL)
		resource := fmt.Sprintf("https://%s", u.Host)

		// authorizer, err := auth.NewAuthorizerFromEnvironmentWithResource(resource)
		authorizer, err := c.newAuthorizerWithEnvVars(auth.NewAuthorizerFromEnvironmentWithResource, resource, c.Env)
		if err != nil {
			return "", 0, err
		}
		c.authorizer = authorizer
	}

	return c.getToken()
}

// GetSiteURL gets SharePoint siteURL
func (c *AuthCnfg) GetSiteURL() string { return c.SiteURL }

// GetStrategy gets auth strategy name
func (c *AuthCnfg) GetStrategy() string { return "azureenv" }

// SetAuth authenticates request
// noinspection GoUnusedParameter
func (c *AuthCnfg) SetAuth(req *http.Request, httpClient *gosip.SPClient) error {
	if _, _, err := c.GetAuth(); err != nil {
		return err
	}
	_, err := c.authorizer.WithAuthorization()(preparer{}).Prepare(req)
	return err
}

// newAuthorizerWithEnvVars sets environment variables and unset them after authorizerFactory code read them
func (c *AuthCnfg) newAuthorizerWithEnvVars(
	authorizerFactory func(resourceBaseURI string) (autorest.Authorizer, error),
	resourceBaseURI string,
	envVars map[string]string,
) (autorest.Authorizer, error) {

	// Set environment variables
	curEnvVars := map[string]string{}
	for key, val := range c.Env {
		if curVal, ok := os.LookupEnv(key); ok {
			curEnvVars[key] = curVal
		}
		os.Setenv(key, val)
	}

	// Get authorizer
	authorizer, err := authorizerFactory(resourceBaseURI)

	// Unset environment variables
	for key := range c.Env {
		prevVal, ok := curEnvVars[key]
		if ok {
			os.Setenv(key, prevVal)
		} else {
			os.Unsetenv(key)
		}
	}

	return authorizer, err
}

// Getting token with prepare for external usage scenarious
func (c *AuthCnfg) getToken() (string, int64, error) {
	req, _ := http.NewRequest("GET", c.SiteURL, nil)
	req, err := c.authorizer.WithAuthorization()(preparer{}).Prepare(req)
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
	return token, j.Exp, nil
}

// Preparer implements autorest.Preparer interface
type preparer struct{}

// Prepare satisfies autorest.Preparer interface
func (p preparer) Prepare(req *http.Request) (*http.Request, error) { return req, nil }
