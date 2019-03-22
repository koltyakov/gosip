package tmg

import (
	"testing"

	h "github.com/koltyakov/gosip/test/helpers"
)

var (
	cnfgPath = "./config/private.tmg.json"
)

func TestGettingAuthToken(t *testing.T) {
	err := h.CheckAuth(
		&AuthCnfg{},
		cnfgPath,
		[]string{"SiteURL", "Username", "Password"},
	)
	if err != nil {
		t.Error(err)
	}
}

func TestGettingDigest(t *testing.T) {
	err := h.CheckDigest(&AuthCnfg{}, cnfgPath)
	if err != nil {
		t.Error(err)
	}
}