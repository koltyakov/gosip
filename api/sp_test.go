package api

import (
	"fmt"
	"os"
	"path"
	"runtime"
	"testing"

	"github.com/koltyakov/gosip"
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
	}

	if envCode != "" && envResolver[envCode] != nil {
		spClient = envResolver[envCode]()
	}

	setHeadersPresets()
}

func TestSP(t *testing.T) {
	checkClient(t)

	t.Run("ToURL", func(t *testing.T) {
		sp := NewSP(spClient)
		if sp.ToURL() != spClient.AuthCnfg.GetSiteURL() {
			t.Errorf(
				"incorrect site URL, expected \"%s\", received \"%s\"",
				spClient.AuthCnfg.GetSiteURL(),
				sp.ToURL(),
			)
		}
	})

	t.Run("Conf", func(t *testing.T) {
		sp := NewSP(spClient)
		sp.config = nil
		sp.Conf(headers.verbose)
		if sp.config != headers.verbose {
			t.Errorf("failed to apply config")
		}
	})

	t.Run("Web", func(t *testing.T) {
		sp := NewSP(spClient)
		if sp.Web() == nil {
			t.Errorf("failed to get Web object")
		}
	})
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
