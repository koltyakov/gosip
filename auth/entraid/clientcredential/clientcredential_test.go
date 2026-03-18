package clientcredential

import (
	"context"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/koltyakov/gosip/auth/entraid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var cnfgPath = "./config/private.spo-entraid-clientcredential.json"

// allow a mock authProvider to be injected during tests
func (c *AuthCnfg) setCredential(cred azcore.TokenCredential) { c.AuthProvider = cred }

type MockCredential struct {
	mock.Mock
}

func newMockCredential() *MockCredential { return &MockCredential{} }

func (m *MockCredential) GetToken(ctx context.Context, opts policy.TokenRequestOptions) (azcore.AccessToken, error) {
	args := m.Called(ctx, opts)
	return args.Get(0).(azcore.AccessToken), args.Error(1)
}

func TestParseConfig_InvalidJSON(t *testing.T) {
	auth := &AuthCnfg{}
	err := auth.ParseConfig([]byte(`not json at all`))
	assert.Error(t, err)
}

func TestReadConfig_FileNotFound(t *testing.T) {
	auth := &AuthCnfg{}
	err := auth.ReadConfig("./config/nonexistent.json")
	assert.Error(t, err)
}

func TestParseConfig(t *testing.T) {
	auth := &AuthCnfg{}
	err := auth.ParseConfig([]byte(`{"siteUrl":"https://contoso.sharepoint.com", "tenantId":"contoso", "clientId":"mocked-client", "certPath":"cert.pfx"}`))
	assert.NoError(t, err)
	assert.Equal(t, "https://contoso.sharepoint.com", auth.SiteURL)
}

func TestGetAuth_ReturnsAToken(t *testing.T) {
	entraid.TokenCache.Flush()
	mockCred := newMockCredential()
	mockCred.On("GetToken", mock.Anything, mock.Anything).
		Return(azcore.AccessToken{Token: "mocked-token", ExpiresOn: time.Now().Add(1 * time.Hour)}, nil)

	auth := &AuthCnfg{Base: entraid.Base{SiteURL: "https://contoso.sharepoint.com", TenantID: "contoso", ClientID: "mocked-client"}}
	auth.setCredential(mockCred)

	token, exp, err := auth.GetAuth()
	assert.NoError(t, err)
	assert.Equal(t, "mocked-token", token)
	assert.True(t, exp > 0)
}

func TestGetAuth_UsesCachedToken(t *testing.T) {
	entraid.TokenCache.Flush()
	mockCred := newMockCredential()
	mockCred.On("GetToken", mock.Anything, mock.Anything).
		Return(azcore.AccessToken{Token: "cached-token", ExpiresOn: time.Now().Add(1 * time.Hour)}, nil).Once()

	auth := &AuthCnfg{Base: entraid.Base{SiteURL: "https://contoso.sharepoint.com", TenantID: "contoso", ClientID: "mocked-client"}}
	auth.setCredential(mockCred)

	token, exp, err := auth.GetAuth()
	assert.NoError(t, err)
	assert.Equal(t, "cached-token", token)
	assert.True(t, exp > 0)

	// second call should hit cache, not the credential
	token, exp, err = auth.GetAuth()
	assert.NoError(t, err)
	assert.Equal(t, "cached-token", token)
	assert.True(t, exp > 0)

	mockCred.AssertNumberOfCalls(t, "GetToken", 1)
	mockCred.AssertExpectations(t)
}

func TestGetAuth_FetchesNewTokenAfterCacheExpiry(t *testing.T) {
	entraid.TokenCache.Flush()

	imminentExpiryTime := time.Now().Add(61 * time.Second)
	freshTokenExpiryTime := time.Now().Add(1 * time.Hour)

	mockCred := newMockCredential()
	// inject a token expiring imminently
	mockCred.On("GetToken", mock.Anything, mock.Anything).
		Return(azcore.AccessToken{Token: "expiring-token", ExpiresOn: imminentExpiryTime}, nil).Once()
	mockCred.On("GetToken", mock.Anything, mock.Anything).
		Return(azcore.AccessToken{Token: "fresh-token", ExpiresOn: freshTokenExpiryTime}, nil).Once()

	auth := &AuthCnfg{Base: entraid.Base{SiteURL: "https://contoso.sharepoint.com", TenantID: "contoso", ClientID: "mocked-client"}}
	auth.setCredential(mockCred)

	token, exp, err := auth.GetAuth()
	assert.NoError(t, err)
	assert.Equal(t, "expiring-token", token)
	assert.Equal(t, imminentExpiryTime.Add(-60*time.Second).Unix(), exp)

	// cached tokens have a TTL of TokenExpiry - 60 seconds
	// thus "expiring-token" token is cached for 1 seconds
	time.Sleep(2 * time.Second) // wait until "expiring-token" expires in cache

	token, exp, err = auth.GetAuth() // fetch and cache a new token
	assert.NoError(t, err)
	assert.Equal(t, "fresh-token", token)
	assert.Equal(t, freshTokenExpiryTime.Add(-60*time.Second).Unix(), exp)

	mockCred.AssertNumberOfCalls(t, "GetToken", 2)
	mockCred.AssertExpectations(t)
}

func TestParseConfig_AbsoluteCertPath(t *testing.T) {
	var absPath, jsonPath string
	if runtime.GOOS == "windows" {
		absPath = `C:\certs\cert.pfx`
		jsonPath = `C:\\certs\\cert.pfx`
	} else {
		absPath = "/etc/certs/cert.pfx"
		jsonPath = absPath
	}

	auth := &AuthCnfg{privateFile: "/some/dir/private.json"}
	err := auth.ParseConfig([]byte(`{"siteUrl":"https://contoso.sharepoint.com","tenantId":"t","clientId":"c","certPath":"` + jsonPath + `"}`))
	assert.NoError(t, err)
	assert.Equal(t, absPath, auth.CertPath, "absolute CertPath should be preserved as-is")
}

func TestParseConfig_RelativeCertPath(t *testing.T) {
	auth := &AuthCnfg{privateFile: filepath.Join("some", "dir", "private.json")}
	err := auth.ParseConfig([]byte(`{"siteUrl":"https://contoso.sharepoint.com","tenantId":"t","clientId":"c","certPath":"cert.pfx"}`))
	assert.NoError(t, err)
	assert.Equal(t, filepath.Join("some", "dir", "cert.pfx"), auth.CertPath, "relative CertPath should be joined with config dir")
}

func TestCacheKey_String(t *testing.T) {
	k := entraid.CacheKey{Host: "contoso.sharepoint.com", Strategy: "clientcredential", TenantID: "tenant", ClientID: "client"}
	assert.Equal(t, "contoso.sharepoint.com@clientcredential@tenant@client", k.String())
}

func TestIntegration_FullAuthFlow(t *testing.T) {
	if _, err := os.Stat(cnfgPath); os.IsNotExist(err) {
		t.Skipf("skipping integration test: %s not found", cnfgPath)
	}

	auth := &AuthCnfg{}
	err := auth.ReadConfig(cnfgPath)
	if err != nil {
		t.Fatalf("ReadConfig failed: %v", err)
	}
	assert.NotEmpty(t, auth.SiteURL, "SiteURL should be populated")
	assert.NotEmpty(t, auth.TenantID, "TenantID should be populated")
	assert.NotEmpty(t, auth.ClientID, "ClientID should be populated")

	token, exp, err := auth.GetAuth()
	if err != nil {
		t.Fatalf("GetAuth failed: %v", err)
	}
	assert.NotEmpty(t, token, "token should not be empty")
	assert.True(t, exp > 0, "expiry should be positive")

	req, _ := http.NewRequest("GET", auth.SiteURL, nil)
	err = auth.SetAuth(req, nil)
	assert.NoError(t, err)
	assert.Contains(t, req.Header.Get("Authorization"), "Bearer ")
}
