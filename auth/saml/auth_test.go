package saml

import (
	"flag"
	"os"
	"testing"

	h "github.com/koltyakov/gosip/test/helpers"
	u "github.com/koltyakov/gosip/test/utils"
)

var (
	cnfgPaths = []string{
		"./config/private.spo-user.json",
		"./config/private.spo-adfs.json",
	}
	ci bool
)

func init() {
	flag.BoolVar(&ci, "ci", false, "Continues integration mode")
	flag.Parse()

	if ci { // In CI mode
		cnfgPath := "./config/private.spo-user.ci.json"
		auth := &AuthCnfg{
			SiteURL:  os.Getenv("SPAUTH_SITEURL"),
			Username: os.Getenv("SPAUTH_USERNAME"),
			Password: os.Getenv("SPAUTH_PASSWORD"),
		}
		auth.WriteConfig(u.ResolveCnfgPath(cnfgPath))
		cnfgPaths = []string{cnfgPath}
	}
}

func TestGettingAuthToken(t *testing.T) {
	configsChecked := 0
	for _, cnfgPath := range cnfgPaths {
		if !h.ConfigExists(cnfgPath) {
			continue
		}
		configsChecked++
		err := h.CheckAuth(
			&AuthCnfg{},
			cnfgPath,
			[]string{"SiteURL", "Username", "Password"},
		)
		if err != nil {
			t.Error(err)
		}
	}
	if configsChecked == 0 {
		t.Skip("No auth config(s) provided")
	}
}

func TestGettingDigest(t *testing.T) {
	configsChecked := 0
	for _, cnfgPath := range cnfgPaths {
		if !h.ConfigExists(cnfgPath) {
			continue
		}
		configsChecked++
		err := h.CheckDigest(&AuthCnfg{}, cnfgPath)
		if err != nil {
			t.Error(err)
		}
	}
	if configsChecked == 0 {
		t.Skip("No auth config(s) provided")
	}
}
