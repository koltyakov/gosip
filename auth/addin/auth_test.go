package addin

import (
	"os"
	"testing"

	h "github.com/koltyakov/gosip/test/helpers"
	u "github.com/koltyakov/gosip/test/utils"
)

var (
	cnfgPath = "./config/private.spo-addin.json"
	ci       bool
)

func init() {
	ci = os.Getenv("SPAUTH_CI") == "true"

	if ci { // In CI mode
		cnfgPath = "./config/private.spo-addin.ci.json"
		auth := &AuthCnfg{
			SiteURL:      os.Getenv("SPAUTH_SITEURL"),
			ClientID:     os.Getenv("SPAUTH_CLIENTID"),
			ClientSecret: os.Getenv("SPAUTH_CLIENTSECRET"),
		}
		auth.WriteConfig(u.ResolveCnfgPath(cnfgPath))
	}
}

func TestGettingAuthToken(t *testing.T) {
	if !h.ConfigExists(cnfgPath) {
		t.Skip("No auth config provided")
	}
	err := h.CheckAuth(
		&AuthCnfg{},
		cnfgPath,
		[]string{"SiteURL", "ClientID", "ClientSecret"},
	)
	if err != nil {
		t.Error(err)
	}
}

func TestGettingDigest(t *testing.T) {
	if !h.ConfigExists(cnfgPath) {
		t.Skip("No auth config provided")
	}
	err := h.CheckDigest(&AuthCnfg{}, cnfgPath)
	if err != nil {
		t.Error(err)
	}
}
