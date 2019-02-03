package tmg

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
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

	cacheKey := parsedURL.Host + "@tmg@" + creds.Username + "@" + creds.Password
	// if accessToken, found := storage.Get(cacheKey); found {
	// 	return accessToken.(string), nil
	// }

	endpoint := fmt.Sprintf("%s://%s/CookieAuth.dll?Logon", parsedURL.Scheme, parsedURL.Host)

	params := url.Values{}
	params.Set("curl", "Z2F") // curl=Z2F&reason=0&formdir=2
	params.Set("reason", "0")
	params.Set("formdir", "2") // TODO: get these params automatically
	params.Set("username", creds.Username)
	params.Set("password", creds.Password)

	// TODO: keepalive agent for https

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	resp, err := client.Post(endpoint, "application/x-www-form-urlencoded", strings.NewReader(params.Encode()))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// fmt.Println(resp.StatusCode)
	authCookie := resp.Header.Get("Set-Cookie") // TODO: parse TMG cookie only (?)

	// TODO: ttl detection
	storage.Set(cacheKey, authCookie, 60*60*time.Second)

	return authCookie, nil
}
