// Package clientcredential provides Entra ID certificate authentication for SharePoint Online.
package clientcredential

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/koltyakov/gosip"
	"github.com/koltyakov/gosip/auth/entraid"
	"github.com/koltyakov/gosip/cpass"
)

// AuthCnfg implements [gosip.AuthCnfg] for certificate-based authentication.
type AuthCnfg struct {
	entraid.Base

	CertPath string `json:"certPath"`
	CertPass string `json:"certPass"`

	privateFile string
}

// ReadConfig reads private config with auth options
func (c *AuthCnfg) ReadConfig(privateFile string) error {
	c.privateFile = privateFile
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
	if err := json.Unmarshal(byteValue, &c); err != nil {
		return err
	}
	if c.privateFile != "" && !filepath.IsAbs(c.CertPath) {
		c.CertPath = filepath.Join(filepath.Dir(c.privateFile), c.CertPath)
	}
	crypt := cpass.Cpass(c.MasterKey)
	secret, err := crypt.Decode(c.CertPass)
	if err == nil {
		c.CertPass = secret
	}
	return nil
}

// WriteConfig writes private config with auth options
func (c *AuthCnfg) WriteConfig(privateFile string) error {
	crypt := cpass.Cpass(c.MasterKey)
	secret, err := crypt.Encode(c.CertPass)
	if err != nil {
		return err
	}
	config := &AuthCnfg{
		Base:     entraid.Base{SiteURL: c.SiteURL, TenantID: c.TenantID, ClientID: c.ClientID},
		CertPath: c.CertPath,
		CertPass: secret,
	}
	file, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(privateFile, file, 0644)
}

// initAuthProvider creates and configures the Azure authentication provider.
func (c *AuthCnfg) initAuthProvider() error {
	certData, err := os.ReadFile(c.CertPath)
	if err != nil {
		return fmt.Errorf("failed to read certificate: %w", err)
	}
	certs, key, err := azidentity.ParseCertificates(certData, []byte(c.CertPass))
	if err != nil {
		return fmt.Errorf("failed to parse certificate: %w", err)
	}
	opts := &azidentity.ClientCertificateCredentialOptions{
		ClientOptions: azcore.ClientOptions{
			Cloud: entraid.CloudFromSiteURL(c.SiteURL),
			Telemetry: policy.TelemetryOptions{
				Disabled: true,
			},
		},
	}
	cred, err := azidentity.NewClientCertificateCredential(c.TenantID, c.ClientID, certs, key, opts)
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
	return c.GetToken(c.GetStrategy())
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
func (c *AuthCnfg) GetStrategy() string { return "clientcredential" }
