package ntlm

import (
	"net/url"
	"time"

	cache "github.com/patrickmn/go-cache"
)

var (
	storage = cache.New(5*time.Minute, 10*time.Minute)
)

// GetAuth : get auth
func GetAuth(c *AuthCnfg) (string, error) {
	parsedURL, err := url.Parse(c.SiteURL)
	if err != nil {
		return "", err
	}

	cacheKey := parsedURL.Host + "@ntlm@" + c.Username + "@" + c.Password
	if accessToken, found := storage.Get(cacheKey); found {
		return accessToken.(string), nil
	}

	return "", nil
}
