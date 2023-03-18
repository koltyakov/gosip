package tmg

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/patrickmn/go-cache"
)

var (
	storage = cache.New(5*time.Minute, 10*time.Minute)
)

// GetAuth gets authentication
func GetAuth(c *AuthCnfg) (string, int64, error) {
	if c.client == nil {
		c.client = &http.Client{}
	}

	parsedURL, err := url.Parse(c.SiteURL)
	if err != nil {
		return "", 0, err
	}

	cacheKey := parsedURL.Host + "@" + c.GetStrategy() + "@" + c.Username + "@" + c.Password
	if accessToken, exp, found := storage.GetWithExpiration(cacheKey); found {
		return accessToken.(string), exp.Unix(), nil
	}

	redirect, err := detectCookieAuthURL(c, c.SiteURL)
	if err != nil {
		return "", 0, err
	}

	endpoint := fmt.Sprintf("%s://%s/CookieAuth.dll?Logon", parsedURL.Scheme, parsedURL.Host)

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

	// client := &http.Client{
	// 	CheckRedirect: func(req *http.Request, via []*http.Request) error {
	// 		return http.ErrUseLastResponse
	// 	},
	// }
	c.client.CheckRedirect = doNotCheckRedirect

	resp, err := c.client.Post(endpoint, "application/x-www-form-urlencoded", strings.NewReader(params.Encode()))
	if err != nil {
		return "", 0, err
	}
	defer func() {
		if resp != nil && resp.Body != nil {
			_ = resp.Body.Close()
		}
	}()

	if _, err := io.Copy(io.Discard, resp.Body); err != nil {
		return "", 0, err
	}

	// fmt.Println(resp.StatusCode)
	authCookie := resp.Header.Get("Set-Cookie") // TODO: parse TMG cookie only (?)

	// TODO: ttl detection
	expiry := time.Hour
	exp := time.Now().Add(expiry).Unix()
	storage.Set(cacheKey, authCookie, expiry)

	return authCookie, exp, nil
}

func detectCookieAuthURL(c *AuthCnfg, siteURL string) (*url.URL, error) {
	if c.client == nil {
		c.client = &http.Client{}
	}

	// client := &http.Client{
	// 	CheckRedirect: func(req *http.Request, via []*http.Request) error {
	// 		return http.ErrUseLastResponse
	// 	},
	// }
	c.client.CheckRedirect = doNotCheckRedirect

	req, err := http.NewRequest("GET", siteURL, nil)
	if err != nil {
		return nil, err
	}

	// req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8")
	// req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/71.0.3578.98 Safari/537.36")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		if resp != nil && resp.Body != nil {
			_ = resp.Body.Close()
		}
	}()

	if _, err := io.Copy(io.Discard, resp.Body); err != nil {
		return nil, err
	}

	redirect, err := resp.Location()
	if err != nil {
		return nil, err
	}

	return redirect, nil
}

// doNotCheckRedirect *http.Client CheckRedirect callback to ignore redirects
func doNotCheckRedirect(_ *http.Request, _ []*http.Request) error {
	return http.ErrUseLastResponse
}
