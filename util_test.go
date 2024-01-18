package gosip

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
)

type AnonymousCnfg struct {
	SiteURL  string `json:"siteUrl"` // SPSite or SPWeb URL, which is the context target for the API calls
	Strategy string
}

// ReadConfig reads private config with auth options
// noinspection GoUnusedParameter
func (c *AnonymousCnfg) ReadConfig(privateFile string) error { return nil }

// ParseConfig parses credentials from a provided JSON byte array content
// noinspection GoUnusedParameter
func (c *AnonymousCnfg) ParseConfig(bytesValue []byte) error { return nil }

// WriteConfig : writes private config with auth options
// noinspection GoUnusedParameter
func (c *AnonymousCnfg) WriteConfig(privateFile string) error { return nil }

// GetAuth : authenticates, receives access token
func (c *AnonymousCnfg) GetAuth(_ context.Context) (string, int64, error) { return "", 0, nil }

// GetSiteURL : gets siteURL
func (c *AnonymousCnfg) GetSiteURL() string { return c.SiteURL }

// GetStrategy : gets auth strategy name
func (c *AnonymousCnfg) GetStrategy() string {
	if c.Strategy != "" {
		return c.Strategy
	}
	return "anonymous"
}

// SetAuth : authenticate request
// noinspection GoUnusedParameter
func (c *AnonymousCnfg) SetAuth(req *http.Request, httpClient *SPClient) error {
	if c.SiteURL == "http://restricted" {
		return fmt.Errorf("can't set auth")
	}
	return nil
}

// Fake server bootstrap helper

func startFakeServer(addr string, handler http.Handler) (io.Closer, error) {
	srv := &http.Server{Addr: addr, Handler: handler}
	if addr == "" {
		addr = ":8989"
	}

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}

	go func() {
		_ = srv.Serve(listener.(*net.TCPListener))
	}()

	return listener, nil
}
