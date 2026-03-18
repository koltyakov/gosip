// Package interactivebrowser provides Entra ID interactive browser authentication for SharePoint Online.
package interactivebrowser

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/koltyakov/gosip"
	"github.com/koltyakov/gosip/auth/entraid"
)

// AuthCnfg implements [gosip.AuthCnfg] for interactive browser authentication.
type AuthCnfg struct {
	entraid.Base

	RedirectURL string `json:"redirectUrl"`
}

// ReadConfig reads private config with auth options
func (c *AuthCnfg) ReadConfig(privateFile string) error {
	f, err := os.Open(privateFile)
	if err != nil {
		return err
	}
	defer f.Close()
	byteValue, err := io.ReadAll(f)
	if err != nil {
		return err
	}
	return c.ParseConfig(byteValue)
}

// ParseConfig parses credentials from a provided JSON byte array content
func (c *AuthCnfg) ParseConfig(byteValue []byte) error {
	return json.Unmarshal(byteValue, c)
}

// WriteConfig writes private config with auth options
func (c *AuthCnfg) WriteConfig(privateFile string) error {
	config := &AuthCnfg{
		Base:        entraid.Base{SiteURL: c.SiteURL, TenantID: c.TenantID, ClientID: c.ClientID},
		RedirectURL: c.RedirectURL,
	}
	file, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(privateFile, file, 0644)
}

// initAuthProvider creates and configures the Azure authentication provider.
func (c *AuthCnfg) initAuthProvider() error {
	if c.RedirectURL == "" {
		return fmt.Errorf("redirectUrl is required")
	}

	opts := &azidentity.InteractiveBrowserCredentialOptions{
		ClientOptions: azcore.ClientOptions{
			Cloud: entraid.CloudFromSiteURL(c.SiteURL),
			Telemetry: policy.TelemetryOptions{
				Disabled: true,
			},
		},
		ClientID:    c.ClientID,
		TenantID:    c.TenantID,
		RedirectURL: c.RedirectURL,
	}
	cred, err := azidentity.NewInteractiveBrowserCredential(opts)
	if err != nil {
		return fmt.Errorf("failed to create authProvider: %w", err)
	}
	c.AuthProvider = cred
	return nil
}

// GetAuth returns an access token and its cache expiry as a Unix timestamp.
func (c *AuthCnfg) GetAuth() (string, int64, error) {
	if c.AuthProvider == nil {
		if err := c.initAuthProvider(); err != nil {
			return "", 0, err
		}
	}
	token, exp, err := c.GetToken(c.GetStrategy())
	if err != nil {
		if strings.Contains(err.Error(), "AADSTS7000218") {
			return "", 0, fmt.Errorf("%w\n\nhint: enable \"Allow public client flows\" in your Entra ID app registration, and ensure the platform is \"Mobile and desktop applications\", not \"Web\"", err)
		}
		return "", 0, err
	}
	return token, exp, nil
}

// SetAuth sets the Bearer token on the request.
func (c *AuthCnfg) SetAuth(req *http.Request, httpClient *gosip.SPClient) error {
	token, _, err := c.GetAuth()
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	return nil
}

// GetStrategy returns the strategy identifier.
func (c *AuthCnfg) GetStrategy() string { return "interactivebrowser" }
