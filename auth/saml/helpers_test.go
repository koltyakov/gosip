package saml

import (
	"context"
	"testing"
)

func TestHelpersEdgeCases(t *testing.T) {

	t.Run("GetAuth/EmptySiteURL", func(t *testing.T) {
		cnfg := &AuthCnfg{SiteURL: ""}
		if _, _, err := GetAuth(context.Background(), cnfg); err == nil {
			t.Error("empty SiteURL should not go")
		}
	})

	t.Run("getSecurityTokenWithAdfs", func(t *testing.T) {
		cnfg := &AuthCnfg{SiteURL: ""}
		if _, _, err := getSecurityTokenWithAdfs(context.Background(), "wrong", cnfg); err == nil {
			t.Error("wrong adfsURL should not go")
		}
		if _, _, err := getSecurityTokenWithAdfs(context.Background(), "http://wrong", cnfg); err == nil {
			t.Error("wrong adfsURL should not go")
		}
	})

}
