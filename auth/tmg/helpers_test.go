package tmg

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

}
