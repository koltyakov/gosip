package entraid

import (
	"context"
	"fmt"
	"net/url"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
)

// Base provides shared functionality for Entra ID auth strategies.
type Base struct {
	SiteURL  string `json:"siteUrl"`
	TenantID string `json:"tenantId"`
	ClientID string `json:"clientId"`

	AuthProvider azcore.TokenCredential
	MasterKey    string
}

// GetToken returns a cached access token or fetches a new one.
// The strategy parameter is used as part of the cache key.
func (b *Base) GetToken(strategy string) (string, int64, error) {
	parsedURL, err := url.Parse(b.SiteURL)
	if err != nil {
		return "", 0, err
	}
	key := CacheKey{Host: parsedURL.Host, Strategy: strategy, TenantID: b.TenantID, ClientID: b.ClientID}
	if accessToken, exp, found := GetCachedToken(key.String()); found {
		return accessToken, exp, nil
	}

	token, err := b.AuthProvider.GetToken(context.Background(), policy.TokenRequestOptions{
		Scopes: []string{fmt.Sprintf("https://%s/.default", parsedURL.Host)},
	})
	if err != nil {
		return "", 0, err
	}

	exp := SetCachedToken(key.String(), token.Token, token.ExpiresOn)
	return token.Token, exp, nil
}

// GetSiteURL returns the configured SharePoint site URL.
func (b *Base) GetSiteURL() string { return b.SiteURL }

// SetMasterkey sets the master key used for credential encryption.
func (b *Base) SetMasterkey(masterKey string) { b.MasterKey = masterKey }
