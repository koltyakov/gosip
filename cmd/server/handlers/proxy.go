package handlers

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/koltyakov/gosip"
)

// Proxy : ...
func Proxy(ctx gosip.AuthCnfg) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		client := &gosip.SPClient{
			AuthCnfg: ctx,
		}

		endpoint := ctx.GetSiteURL() + r.RequestURI

		req, err := http.NewRequest(r.Method, endpoint, r.Body)
		if err != nil {
			message := fmt.Sprintf("unable to create a request: %v", err)
			http.Error(w, message, http.StatusBadRequest)
			return
		}

		for name, headers := range r.Header {
			name = strings.ToLower(name)
			for _, h := range headers {
				req.Header.Add(name, h)
			}
		}

		fmt.Printf("requesting endpoint: %s\n", endpoint)
		resp, err := client.Execute(req)
		if err != nil {
			message := fmt.Sprintf("unable to request: %v\n", err)
			http.Error(w, message, http.StatusBadRequest)
			return
		}
		defer resp.Body.Close()

		for name, headers := range resp.Header {
			name = strings.ToLower(name)
			for _, h := range headers {
				w.Header().Add(name, h)
			}
		}

		w.WriteHeader(resp.StatusCode)

		io.Copy(w, resp.Body)
	}
}
