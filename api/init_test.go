package api

import (
	"fmt"
	"os"
	"path"
	"runtime"
	"testing"

	"github.com/koltyakov/gosip"
	ntlm "github.com/koltyakov/gosip/auth/ntlm"
	saml "github.com/koltyakov/gosip/auth/saml"
	h "github.com/koltyakov/gosip/test/helpers"
)

var (
	ci         bool
	heavyTests bool
	envCode    string
	spClient   *gosip.SPClient
	headers    struct {
		verbose         *RequestConfig
		minimalmetadata *RequestConfig
		nometadata      *RequestConfig
	}
)

func init() {
	ci = os.Getenv("SPAUTH_CI") == "true"
	heavyTests = os.Getenv("SPAPI_HEAVY_TESTS") == "true"
	envCode = os.Getenv("SPAUTH_ENVCODE")

	if envCode == "" && !ci {
		envCode = "spo"
	}

	envResolver := map[string]func() *gosip.SPClient{
		"spo": func() *gosip.SPClient {
			cnfgPath := "./config/integration/private.spo.json"
			auth := &saml.AuthCnfg{}
			if ci {
				auth.SiteURL = os.Getenv("ENV_SPO_SITEURL")
				auth.Username = os.Getenv("ENV_SPO_USERNAME")
				auth.Password = os.Getenv("ENV_SPO_PASSWORD")
			} else {
				if err := auth.ReadConfig(resolveCnfgPath(cnfgPath)); err != nil {
					return nil
				}
			}
			if err := h.CheckAuthProps(auth, []string{"SiteURL", "Username", "Password"}); err != nil {
				return nil
			}
			client := &gosip.SPClient{AuthCnfg: auth}
			// Pre-auth for tests not include auth timing involved
			if _, err := client.AuthCnfg.GetAuth(); err != nil {
				fmt.Printf("can't auth, %s\n", err)
				// Force all test being skipped in case of auth issues
				return nil
			}
			return client
		},
		"2013": func() *gosip.SPClient {
			cnfgPath := "./config/integration/private.2013.json"
			auth := &ntlm.AuthCnfg{}
			if err := auth.ReadConfig(resolveCnfgPath(cnfgPath)); err != nil {
				return nil
			}
			if err := h.CheckAuthProps(auth, []string{"SiteURL", "Username", "Password", "Domain"}); err != nil {
				return nil
			}
			client := &gosip.SPClient{AuthCnfg: auth}
			// Pre-auth for tests not include auth timing involved
			if _, err := client.AuthCnfg.GetAuth(); err != nil {
				fmt.Printf("can't auth, %s\n", err)
				// Force all test being skipped in case of auth issues
				return nil
			}
			return client
		},
	}

	if envCode != "" && envResolver[envCode] != nil {
		spClient = envResolver[envCode]()
	}

	setHeadersPresets()
}

func resolveCnfgPath(relativePath string) string {
	_, filename, _, _ := runtime.Caller(1)
	fmt.Println(filename)
	return path.Join(path.Dir(filename), "..", relativePath)
}

func checkClient(t *testing.T) {
	if spClient == nil {
		t.Skip("no auth context provided")
	}
}

func setHeadersPresets() {
	headers.verbose = HeadersPresets.Verbose
	headers.minimalmetadata = HeadersPresets.Minimalmetadata
	headers.nometadata = HeadersPresets.Nometadata
}
