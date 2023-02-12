// Package device implements AAD Device Auth Flow
// See more: https://docs.microsoft.com/en-us/azure/developer/go/azure-sdk-authorization#use-device-token-authentication
// Amongst supported platform versions are:
//   - SharePoint Online + Azure
package device

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"

	"github.com/Azure/go-autorest/autorest/adal"
	"github.com/Azure/go-autorest/autorest/azure/auth"
	"github.com/koltyakov/gosip"
	"github.com/koltyakov/gosip/cpass"
)

var (
	tokenCache = map[string]*adal.ServicePrincipalToken{} // ToDo: Replace with sync.Map
	crypter    = cpass.Cpass("")
)

// AuthCnfg - AAD Device Flow auth config structure
/* Config sample:
{
  "siteUrl": "https://contoso.sharepoint.com/sites/test",
	"clientId": "61367a97-562c-4372-a9ee-b35307abdd26",
	"tenantId": "3f83fe32-29b2-488e-8c3f-c8b7a2e19a2f"
}
*/
type AuthCnfg struct {
	SiteURL  string `json:"siteUrl"`  // SPSite or SPWeb URL, which is the context target for the API calls
	ClientID string `json:"clientId"` // Azure AD App Registration Client ID
	TenantID string `json:"tenantId"` // Azure AD App Registration Tenant ID
}

// ReadConfig reads private config with auth options
func (c *AuthCnfg) ReadConfig(privateFile string) error {
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
	return json.Unmarshal(byteValue, &c)
}

// WriteConfig writes private config with auth options
func (c *AuthCnfg) WriteConfig(privateFile string) error {
	config := &AuthCnfg{
		SiteURL:  c.SiteURL,
		ClientID: c.ClientID,
		TenantID: c.TenantID,
	}
	file, _ := json.MarshalIndent(config, "", "  ")
	return ioutil.WriteFile(privateFile, file, 0644)
}

// GetAuth authenticates, receives access token
func (c *AuthCnfg) GetAuth() (string, int64, error) {
	u, _ := url.Parse(c.SiteURL)
	resource := fmt.Sprintf("https://%s", u.Host)

	// Check cached token per resource
	token := tokenCache[resource]

	// Check disk cache
	if token == nil {
		token, _ = c.getTokenDiskCache()
	}

	if token != nil {
		// Return cached token if not expired
		if !token.Token().IsExpired() {
			return token.Token().AccessToken, token.Token().Expires().Unix(), nil
		}
		// Expired, try to refresh
		if err := token.Refresh(); err == nil {
			// Cache refreshed token
			_ = c.cacheTokenToDisk(token)
			// Return refreshed token
			return token.Token().AccessToken, token.Token().Expires().Unix(), nil
		}
		// Failed to refresh, initiating for the device auth flow
	}

	config := auth.NewDeviceFlowConfig(c.ClientID, c.TenantID)
	config.Resource = resource

	token, err := config.ServicePrincipalToken()
	if err != nil {
		return "", 0, err
	}

	_ = c.cacheTokenToDisk(token)

	tokenCache[resource] = token
	return token.Token().AccessToken, token.Token().Expires().Unix(), nil
}

// GetSiteURL gets SharePoint siteURL
func (c *AuthCnfg) GetSiteURL() string { return c.SiteURL }

// GetStrategy gets auth strategy name
func (c *AuthCnfg) GetStrategy() string { return "device" }

// SetAuth authenticates request
// noinspection GoUnusedParameter
func (c *AuthCnfg) SetAuth(req *http.Request, httpClient *gosip.SPClient) error {
	accessToken, _, err := c.GetAuth()
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)
	return nil
}

// === File system token caching helpers === //

// CleanTokenCache removes token information
func (c *AuthCnfg) CleanTokenCache() error {
	tokenCachePath := c.getTokenCachePath()

	delete(tokenCache, c.ClientID)
	if err := os.Remove(tokenCachePath); err != nil {
		return err
	}
	return nil
}

// cacheTokenToDisk writes serialized token to temporary cache file
func (c *AuthCnfg) cacheTokenToDisk(token *adal.ServicePrincipalToken) error {
	tmpDir := filepath.Join(os.TempDir(), "gosip")
	tokenCachePath := c.getTokenCachePath()

	tokenCache, err := token.MarshalJSON()
	if err != nil {
		return err
	}
	tokenCacheE, _ := crypter.Encode(fmt.Sprintf("%s", tokenCache))
	tokenCache = []byte(tokenCacheE)

	_ = os.MkdirAll(tmpDir, os.ModePerm)
	if err := ioutil.WriteFile(tokenCachePath, tokenCache, 0644); err != nil {
		return err
	}
	return nil
}

// getTokenDiskCache reads token from temporary cache file
func (c *AuthCnfg) getTokenDiskCache() (*adal.ServicePrincipalToken, error) {
	tokenCachePath := c.getTokenCachePath()

	tokenCache, err := ioutil.ReadFile(tokenCachePath)
	if err != nil {
		return nil, err
	}
	tokenCacheD, _ := crypter.Decode(fmt.Sprintf("%s", tokenCache))
	tokenCache = []byte(tokenCacheD)

	token := &adal.ServicePrincipalToken{}
	if err := token.UnmarshalJSON(tokenCache); err != nil {
		return nil, err
	}
	return token, nil
}

// getTokenCachePath gets local file system file path with token cache
func (c *AuthCnfg) getTokenCachePath() string {
	tmpDir := filepath.Join(os.TempDir(), "gosip")
	return filepath.Join(tmpDir, c.GetStrategy()+"_"+c.ClientID)
}
