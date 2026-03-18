// Package devicecode provides Entra ID device code authentication for SharePoint Online.
package devicecode

import (
	"context"
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

// AuthCnfg implements [gosip.AuthCnfg] for device code authentication.
type AuthCnfg struct {
	entraid.Base

	// UserPrompt is called with the device code message. If nil, defaults to printing to stdout.
	UserPrompt func(ctx context.Context, msg azidentity.DeviceCodeMessage) error `json:"-"`
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
		Base: entraid.Base{SiteURL: c.SiteURL, TenantID: c.TenantID, ClientID: c.ClientID},
	}
	file, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(privateFile, file, 0644)
}

// initAuthProvider creates and configures the Azure authentication provider.
func (c *AuthCnfg) initAuthProvider() error {
	userPrompt := c.UserPrompt
	if userPrompt == nil {
		userPrompt = func(_ context.Context, msg azidentity.DeviceCodeMessage) error {
			fmt.Println(msg.Message)
			return nil
		}
	}
	opts := &azidentity.DeviceCodeCredentialOptions{
		ClientOptions: azcore.ClientOptions{
			Cloud: entraid.CloudFromSiteURL(c.SiteURL),
			Telemetry: policy.TelemetryOptions{
				Disabled: true,
			},
		},
		ClientID:   c.ClientID,
		TenantID:   c.TenantID,
		UserPrompt: userPrompt,
	}
	cred, err := azidentity.NewDeviceCodeCredential(opts)
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
			return "", 0, fmt.Errorf("%w\n\nhint: enable \"Allow public client flows\" in your Entra ID app registration settings", err)
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
func (c *AuthCnfg) GetStrategy() string { return "devicecode" }
