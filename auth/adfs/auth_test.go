package adfs

import (
	"testing"

	h "github.com/koltyakov/gosip/test/helpers"
)

var (
	cnfgPaths = []string{
		"./config/private.onprem-adfs.json",
		"./config/private.onprem-wap.json",
		"./config/private.onprem-wap-adfs.json",
	}
)

func TestGettingAuthToken(t *testing.T) {
	for _, cnfgPath := range cnfgPaths {
		if !h.ConfigExists(cnfgPath) {
			continue
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
}

func TestGettingDigest(t *testing.T) {
	for _, cnfgPath := range cnfgPaths {
		if !h.ConfigExists(cnfgPath) {
			continue
		}
		err := h.CheckDigest(&AuthCnfg{}, cnfgPath)
		if err != nil {
			t.Error(err)
		}
	}
}
