package fba

import (
	"testing"

	h "github.com/koltyakov/gosip/test/helpers"
)

var (
	cnfgPath = "./config/private.onprem-fba.json"
)

func TestGettingAuthToken(t *testing.T) {
	if !h.ConfigExists(cnfgPath) {
		t.Skip("No auth config provided")
	}
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
	if !h.ConfigExists(cnfgPath) {
		t.Skip("No auth config provided")
	}
	err := h.CheckDigest(&AuthCnfg{}, cnfgPath)
	if err != nil {
		t.Error(err)
	}
}
