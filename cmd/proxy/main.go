package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/koltyakov/gosip"
	"github.com/koltyakov/gosip/auth/addin"
	"github.com/koltyakov/gosip/auth/adfs"
	"github.com/koltyakov/gosip/auth/fba"
	"github.com/koltyakov/gosip/auth/ntlm"
	"github.com/koltyakov/gosip/auth/saml"
	"github.com/koltyakov/gosip/auth/tmg"
)

var debug bool

func main() {

	var strategy string
	var config string
	var port int
	var sslKey string
	var sslCert string

	flag.StringVar(&strategy, "strategy", "", "Auth strategy")
	flag.StringVar(&config, "config", "", "Config path")
	flag.IntVar(&port, "port", 9090, "Proxy port")
	flag.BoolVar(&debug, "debug", false, "Debug mode")
	flag.StringVar(&sslKey, "sslKey", "", "SSL Key file path")   // openssl genrsa -out private.key 2048
	flag.StringVar(&sslCert, "sslCert", "", "SSL Crt file path") // openssl req -new -x509 -sha256 -key private.key -out public.crt -days 3650

	flag.Parse()

	var auth gosip.AuthCnfg

	switch strategy {
	case "addin":
		auth = &addin.AuthCnfg{}
		break
	case "adfs":
		auth = &adfs.AuthCnfg{}
		break
	case "fba":
		auth = &fba.AuthCnfg{}
		break
	case "ntlm":
		auth = &ntlm.AuthCnfg{}
		break
	case "saml":
		auth = &saml.AuthCnfg{}
		break
	case "tmg":
		auth = &tmg.AuthCnfg{}
		break
	default:
		log.Fatalf("can't resolve the strategy: %s", strategy)
	}

	if config == "" {
		log.Fatalf("config path must be provided")
	}

	err := auth.ReadConfig(config)
	if err != nil {
		log.Fatalf("unable to get config: %v", err)
	}

	http.HandleFunc("/", proxyHandler(auth))

	if sslKey != "" && sslCert != "" {
		log.Fatal(http.ListenAndServeTLS(fmt.Sprintf(":%d", port), sslCert, sslKey, nil))
	} else {
		log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
	}

}

func proxyHandler(ctx gosip.AuthCnfg) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		client := &gosip.SPClient{
			AuthCnfg: ctx,
		}

		siteURL, err := url.Parse(ctx.GetSiteURL())
		if err != nil {
			message := fmt.Sprintf("unable to parse site url: %v", err)
			http.Error(w, message, http.StatusBadRequest)
			return
		}

		endpoint := strings.Replace(ctx.GetSiteURL(), siteURL.Path, "", -1) + r.RequestURI
		if strings.Contains(r.RequestURI, siteURL.Path) == false {
			endpoint = ctx.GetSiteURL() + r.RequestURI
		}

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

		if debug {
			fmt.Printf("requesting endpoint: %s\n", endpoint)
		}
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
