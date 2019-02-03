package basic

import (
	"net/url"
	"time"

	cache "github.com/patrickmn/go-cache"
)

var (
	storage = cache.New(5*time.Minute, 10*time.Minute)
)

// GetAuth : get auth
func GetAuth(creds *AuthCnfg) (string, error) {
	parsedURL, err := url.Parse(creds.SiteURL)
	if err != nil {
		return "", err
	}

	cacheKey := parsedURL.Host + "@basic@" + creds.Username + "@" + creds.Password
	if accessToken, found := storage.Get(cacheKey); found {
		return accessToken.(string), nil
	}

	// ntlmssp.Negotiator()
	return "OK", nil
}
