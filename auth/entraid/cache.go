// Package entraid provides shared utilities for Entra ID auth strategies.
package entraid

import (
	"fmt"
	"time"

	"github.com/patrickmn/go-cache"
)

// TokenCache is a shared token cache for Entra ID auth strategies.
var TokenCache = cache.New(5*time.Minute, 10*time.Minute)

// CacheKey identifies a cached token by host, strategy, tenant, and client.
type CacheKey struct {
	Host     string
	Strategy string
	TenantID string
	ClientID string
}

func (k CacheKey) String() string {
	return fmt.Sprintf("%s@%s@%s@%s", k.Host, k.Strategy, k.TenantID, k.ClientID)
}

// GetCachedToken returns a cached access token and its expiry if found.
func GetCachedToken(key string) (token string, exp int64, found bool) {
	if v, expTime, ok := TokenCache.GetWithExpiration(key); ok {
		token, ok := v.(string)
		if !ok {
			return "", 0, false
		}
		return token, expTime.Unix(), true
	}
	return "", 0, false
}

// SetCachedToken stores an access token in the cache with a TTL derived from the token's
// expiry time minus a 60-second buffer.
func SetCachedToken(key string, token string, expiresOn time.Time) int64 {
	exp := expiresOn.Add(-60 * time.Second)
	TokenCache.Set(key, token, time.Until(exp))
	return exp.Unix()
}
