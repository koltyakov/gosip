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
func GetAuth(c *AuthCnfg) (string, error) {
	parsedURL, err := url.Parse(c.SiteURL)
	if err != nil {
		return "", err
	}

	cacheKey := parsedURL.Host + "@tmg@" + c.Username + "@" + c.Password
	if accessToken, found := storage.Get(cacheKey); found {
		return accessToken.(string), nil
	}

	redirect, err := detectCookieAuthURL(c.SiteURL)
	if err != nil {
		return "", err
	}
	// fmt.Printf("Redirect URL: %s\n", redirect)

	endpoint := fmt.Sprintf("%s://%s/CookieAuth.dll?Logon", parsedURL.Scheme, parsedURL.Host)

	// fmt.Printf("Endpoint: %s\n", endpoint)

	params := url.Values{}

	querystr := strings.Replace(redirect.RawQuery, "GetLogon?", "", 1)
	for _, part := range strings.Split(querystr, "&") {
		p := strings.Split(part, "=")
		if len(p) == 2 {
			params.Set(p[0], p[1])
		}
	}

	params.Set("username", c.Username)
	params.Set("password", c.Password)

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
	expirity := time.Hour
	storage.Set(cacheKey, authCookie, expirity)

	return authCookie, nil
}

func detectCookieAuthURL(siteURL string) (*url.URL, error) {
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	req, err := http.NewRequest("GET", siteURL, nil)
	if err != nil {
		return nil, err
	}

	// req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8")
	// req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/71.0.3578.98 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	redirect, err := resp.Location()
	if err != nil {
		return nil, err
	}

	return redirect, nil
}
