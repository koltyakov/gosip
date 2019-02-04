package gosip

import (
	"errors"
	"net/http"
)

// AuthCnfg : abstract auth interface
type AuthCnfg interface {
	ReadConfig(configPath string) error
	GetAuth() (string, error)
	GetSiteURL() string
	// GetStrategy() string
	SetAuth(req *http.Request) error
}

// SPClient : SharePoint HTTP client struct
type SPClient struct {
	http.Client
	AuthCnfg   AuthCnfg
	ConfigPath string
	// Transport  http.RoundTripper
}

// Execute : SharePoint HTTP client Do method
func (c *SPClient) Execute(req *http.Request) (*http.Response, error) {
	// fmt.Println("Injecting auth call...")
	if c.ConfigPath != "" && c.AuthCnfg.GetSiteURL() == "" {
		c.AuthCnfg.ReadConfig(c.ConfigPath)
	}
	if c.AuthCnfg.GetSiteURL() == "" {
		res := &http.Response{
			Status:     "400 Internal Error",
			StatusCode: 400,
			Request:    req,
		}
		return res, errors.New("Client initialization error")
	}
	// authCookie, err := c.AuthCnfg.GetAuth()
	// if err != nil {
	// 	// fmt.Printf("unable to get auth: %v", err)
	// 	res := &http.Response{
	// 		Status:     "401 Access Denied",
	// 		StatusCode: 401,
	// 		Request:    req,
	// 	}
	// 	return res, err
	// }
	// req.Header.Set("Cookie", authCookie)
	err := c.AuthCnfg.SetAuth(req)
	if err != nil {
		// fmt.Printf("unable to get auth: %v", err)
		res := &http.Response{
			Status:     "401 Access Denied",
			StatusCode: 401,
			Request:    req,
		}
		return res, err
	}
	return c.Do(req)
}
