// Package anon provides anonymous "strategy"
// no auth mechanisms are applied to the requests
package anon

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"os"

	"github.com/koltyakov/gosip"
)

// AuthCnfg - anonymous config structure
/* Config sample:
{
  "siteUrl": "https://contoso/sites/test"
}
*/
type AuthCnfg struct {
	SiteURL string `json:"siteUrl"` // SPSite or SPWeb URL, which is the context target for the API calls
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

	return nil
}

// WriteConfig writes private config with auth options
func (c *AuthCnfg) WriteConfig(privateFile string) error {
	config := &AuthCnfg{
		SiteURL: c.SiteURL,
	}
	file, _ := json.MarshalIndent(config, "", "  ")
	return os.WriteFile(privateFile, file, 0644)
}

// GetAuth authenticates, receives access token
func (c *AuthCnfg) GetAuth(ctx context.Context) (string, int64, error) { return "", 0, nil }

// GetSiteURL gets siteURL
func (c *AuthCnfg) GetSiteURL() string { return c.SiteURL }

// GetStrategy gets auth strategy name
func (c *AuthCnfg) GetStrategy() string { return "anonymous" }

// SetAuth : authenticate request
// noinspection GoUnusedParameter
func (c *AuthCnfg) SetAuth(req *http.Request, httpClient *gosip.SPClient) error { return nil }
