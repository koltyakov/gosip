package addin

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

	t.Run("GetAuth/IncorrectSiteURL", func(t *testing.T) {
		cnfg := &AuthCnfg{SiteURL: "https://wrong"}
		if _, err := GetAuth(cnfg); err == nil {
			t.Error("incorrect SiteURL should not go")
		}
	})

	t.Run("getRealm/EmptyRealm", func(t *testing.T) {
		cnfg := &AuthCnfg{Realm: ""}
		if _, err := GetAuth(cnfg); err == nil {
			t.Error("empty Realm should not go")
		}
	})

	t.Run("getRealm/EmptySiteURL", func(t *testing.T) {
		cnfg := &AuthCnfg{Realm: "any", SiteURL: ""}
		if _, err := GetAuth(cnfg); err == nil {
			t.Error("empty SiteURL should not go")
		}
	})

	t.Run("getRealm/IncorrectSiteURL", func(t *testing.T) {
		cnfg := &AuthCnfg{SiteURL: "https://wrong", Realm: "wrong"}
		if _, err := GetAuth(cnfg); err == nil {
			t.Error("incorrect SiteURL should not go")
		}
	})

}
