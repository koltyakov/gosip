package saml

import (
	"testing"
)

func TestHelpersEdgeCases(t *testing.T) {

	t.Run("GetAuth/EmptySiteURL", func(t *testing.T) {
		cnfg := &AuthCnfg{SiteURL: ""}
		if _, err := GetAuth(cnfg); err == nil {
			t.Error("empty SiteURL should not go")
		}
	})

	t.Run("getSecurityTokenWithAdfs", func(t *testing.T) {
		cnfg := &AuthCnfg{SiteURL: ""}
		if _, _, err := getSecurityTokenWithAdfs("wrong", cnfg); err == nil {
			t.Error("wrong adfsURL should not go")
		}
		if _, _, err := getSecurityTokenWithAdfs("http://wrong", cnfg); err == nil {
			t.Error("wrong adfsURL should not go")
		}
	})

}
